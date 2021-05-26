package dao

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	dbsql "database/sql"
	"hcc/piano/lib/mysql"
	"hcc/piano/model"
)

func sendStmt(sql string, params ...interface{}) (dbsql.Result, error) {
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = stmt.Close()
	}()

	result, err := stmt.Exec(params...)

	if err != nil {
		return result, err
	}

	return result, nil
}

func sendQuery(sql string) (*dbsql.Rows, error) {
	result, err := mysql.Db.Query(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func InsertNetworkBillingInfo(infoList *[]model.NetworkBill) error {
	sql := "INSERT INTO `piano`.`network_billing_info` (`group_id`, `date`, `charge_subnet`, `charge_adaptive_ip`) VALUES "

	for _, info := range *infoList {
		sql += fmt.Sprintf("(%d, DATE(NOW()), %d, %d),",
			info.GroupID,
			info.ChargeSubnet,
			info.ChargeAdaptiveIP)
	}

	sql = strings.TrimSuffix(sql, ",") + " AS `new_info` "
	sql += "ON DUPLICATE KEY UPDATE " +
		"`charge_subnet` = `new_info`.`charge_subnet`, " +
		"`charge_adaptive_ip` = `new_info`.`charge_adaptive_ip`;"

	res, err := sendQuery(sql)
	if res != nil {
		_ = res.Close()
	}
	return err
}

func InsertNodeBillingInfo(infoList *[]model.NodeBill) error {
	sql := "INSERT INTO `piano`.`node_billing_info` (`group_id`, `date`, `node_uuid`, `charge_cpu`, `charge_memory`, `charge_nic`) VALUES "

	for _, info := range *infoList {
		sql += fmt.Sprintf("(%d, DATE(NOW()), '%s', %d, %d, %d),",
			info.GroupID,
			info.NodeUUID,
			info.ChargeCPU,
			info.ChargeMEM,
			info.ChargeNIC)
	}

	sql = strings.TrimSuffix(sql, ",") + " AS `new_info` "
	sql += "ON DUPLICATE KEY UPDATE " +
		"`charge_cpu` = `new_info`.`charge_cpu`, " +
		"`charge_memory` = `new_info`.`charge_memory`, " +
		"`charge_nic` = `new_info`.`charge_nic`;"

	res, err := sendQuery(sql)
	if res != nil {
		_ = res.Close()
	}

	return err
}

func InsertServerBillingInfo(infoList *[]model.ServerBill) error {
	sql := "INSERT INTO `piano`.`server_billing_info` (`group_id`, `date`, `server_uuid`, `charge_traffic`) VALUES "

	for _, info := range *infoList {
		sql += fmt.Sprintf("(%d, DATE(NOW()), '%s', %d),",
			info.GroupID,
			info.ServerUUID,
			info.ChargeTraffic)
	}

	sql = strings.TrimSuffix(sql, ",") + " AS `new_info` "
	sql += "ON DUPLICATE KEY UPDATE " +
		"`charge_traffic` = `new_info`.`charge_traffic`;"

	res, err := sendQuery(sql)
	if res != nil {
		_ = res.Close()
	}

	return err
}

func InsertVolumeBillingInfo(infoList *[]model.VolumeBill) error {
	sql := "INSERT INTO `piano`.`volume_billing_info` (`group_id`, `date`, `charge_ssd`, `charge_hdd`) VALUES "

	for _, info := range *infoList {
		sql += fmt.Sprintf("(%d, DATE(NOW()), %d, %d),",
			info.GroupID,
			info.ChargeSSD,
			info.ChargeHDD)
	}

	sql = strings.TrimSuffix(sql, ",") + " AS `new_info` "
	sql += "ON DUPLICATE KEY UPDATE " +
		"`charge_ssd` = `new_info`.`charge_ssd`," +
		"`charge_hdd` = `new_info`.`charge_hdd`;"

	res, err := sendQuery(sql)
	if res != nil {
		_ = res.Close()
	}

	return err
}

func InsertDailyInfo() error {
	sql := "INSERT INTO `piano`.`daily_info` (`date`, `group_id`, `charge_node`, `charge_server`, `charge_network`, `charge_volume`) SELECT	`current_info`.`date`, `current_info`.`group_id`, `current_info`.`charge_node`,	`current_info`.`charge_network`, `current_info`.`charge_server`, `current_info`.`charge_volume` FROM (SELECT `node_charge`.`date` AS `date`, `node_charge`.`group_id` AS `group_id`, `charge_node`, `charge_network`, `charge_server`, `charge_volume` FROM (SELECT `node`.`group_id`, `node`.`date`, (`node`.`charge_cpu` + `node`.`charge_memory` + `node`.`charge_nic`) AS `charge_node` FROM `piano`.`node_billing_info` AS `node` WHERE `node`.`date` = DATE(NOW()) GROUP BY `node`.`group_id`, `node`.`date`) AS `node_charge` LEFT JOIN (SELECT `net`.`group_id`, `net`.`date`, (`net`.`subnet_charge` + `net`.`adaptive_ip_charge`) AS `charge_network` FROM `piano`.`network_billing_info` AS `net` WHERE `net`.`date` = DATE(NOW())) AS `net_charge` ON `net_charge`.`date` = `node_charge`.`date` AND `net_charge`.`group_id` = `node_charge`.`group_id` LEFT JOIN (SELECT `server`.`group_id`, `server`.`date`, SUM(`server`.`network_traffic` * `server`.`traffic_charge_per_kb`) AS `charge_server` FROM `piano`.`server_billing_info` AS `server` WHERE `server`.`date` = DATE(NOW()) GROUP BY `server`.`group_id`, `server`.`date`) AS `server_charge` ON `server_charge`.`date` = `node_charge`.`date` AND `server_charge`.`group_id` = `node_charge`.`group_id` LEFT JOIN (SELECT `volume`.`group_id`, `volume`.`date`, (``volume`.`hdd_charge` + `volume`.`ssd_charge`) AS `charge_volume` FROM `piano`.`volume_billing_info` AS `volume` WHERE `volume`.`date` = DATE(NOW())) AS `volume_charge` ON `volume_charge`.`date` = `node_charge`.`date` AND `volume_charge`.`group_id` = `node_charge`.`group_id`) AS `current_info` ON DUPLICATE KEY UPDATE `daily_info`.`charge_node` = `current_info`.`charge_node`, `daily_info`.`charge_server` = `current_info`.`charge_server`, `daily_info`.`charge_network` = `current_info`.`charge_network`, `daily_info`.`charge_volume` = `current_info`.`charge_volume`;"

	res, err := sendQuery(sql)
	if res != nil {
		_ = res.Close()
	}

	return err
}

func GetBill(groupID int, start, end, billType string, row, page int) (*dbsql.Rows, error) {
	billIdStart := strconv.Itoa(groupID) + start
	billIdEnd := strconv.Itoa(groupID) + end
	billType = strings.ToLower(billType)
	sql := ""

	switch billType {
	case "daily":
		fallthrough
	case "monthly":
		fallthrough
	case "yearly":
		sql = "SELECT * FROM `piano`.`" + billType + "_bill` WHERE `bill_id` BETWEEN " + billIdStart + " AND " +
			billIdEnd + " LIMIT " + strconv.Itoa(row) + " OFFSET " + strconv.Itoa(row*page) + ";"
	default:
		return nil, errors.New("DAO(GetBill) -> Unsupported billing type")
	}

	res, err := sendQuery(sql)

	return res, err
}

func GetBillInfo(groupID int, date, billType, category string) (*dbsql.Rows, error) {
	billType = strings.ToLower(billType)
	category = strings.ToLower(category)
	dateStart := date
	dateEnd, _ := strconv.Atoi(date)

	switch billType {
	case "daily":
		break
	case "monthly":
		dateEnd += 100
		if dateEnd%10000 > 12 {
			dateEnd += 10000
			dateEnd -= 100
		}
	case "yearly":
		dateEnd += 10000
	default:
		return nil, errors.New("DAO(GetBill) -> Unsupported billing type")
	}

	switch category {
	case "network":
		fallthrough
	case "server":
		fallthrough
	case "node":
		fallthrough
	case "volume":
		break
	default:
		return nil, errors.New("DAO(GetBill) -> Unsupported category")
	}

	sql := "SELECT * FROM `piano`.`" + category + "_billing_info` WHERE `group_id`=" + strconv.Itoa(groupID) + " AND `date` BETWEEN DATE(" + dateStart + ") AND DATE(" + strconv.Itoa(dateEnd) + ");"

	res, err := sendQuery(sql)
	if res != nil {
		_ = res.Close()
	}

	return res, err
}
