package dao

import (
	"fmt"
	"strconv"
	"strings"

	//"hcc/piano/lib/logger"
	"hcc/piano/lib/mysql"
	"hcc/piano/model"

	errors "innogrid.com/hcloud-classic/hcc_errors"
)

func sendStmt(sql string, params ...interface{}) (mysql.Result, *errors.HccError) {
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		return nil, errors.NewHccError(errors.PianoInternalOperationFail, "sql Prepare : "+err.Error())
	}

	defer func() {
		_ = stmt.Close()
	}()

	result, err := stmt.Exec(params...)

	if err != nil {
		return result, errors.NewHccError(errors.PianoInternalOperationFail, "stmt Exec : "+err.Error())
	}

	return result, nil
}

func sendQuery(sql string) (*mysql.Rows, *errors.HccError) {
	result, err := mysql.Db.Query(sql)
	if err != nil {
		return nil, errors.NewHccError(errors.PianoInternalOperationFail, "sql Query : "+err.Error())
	}

	return result, nil
}

func InsertNetworkBillingInfo(infoList *[]model.NetworkBill) *errors.HccError {
	sql := "INSERT INTO `piano`.`network_billing_info` (`group_id`, `date`, `subnet_count`, `adaptiveip_count`, `subnet_charge_per_cnt`, `adaptiveip_charge_per_cnt`, `discount_rate`) VALUES "

	for _, info := range *infoList {
		sql += fmt.Sprintf("(%d, DATE(NOW()), %d, %d, %f, %f, %f),",
			info.GroupID,
			info.SubnetCount,
			info.AIPCount,
			info.SubnetChargePerCnt,
			info.AIPChargePerCnt,
			info.DiscountRate)
	}

	sql = strings.TrimSuffix(sql, ",") + " AS `new_info` "
	sql += "ON DUPLICATE KEY UPDATE " +
		"`subnet_count` = `new_info`.`subnet_count`, " +
		"`adaptiveip_count` = `new_info`.`adaptiveip_count`, " +
		"`subnet_charge_per_cnt` = `new_info`.`subnet_charge_per_cnt`, " +
		"`adaptiveip_charge_per_cnt` = `new_info`.`adaptiveip_charge_per_cnt`, " +
		"`discount_rate` = `new_info`.`discount_rate`;"

	res, err := sendQuery(sql)
	res.Close()
	return err
}

func InsertNodeBillingInfo(infoList *[]model.NodeBill) *errors.HccError {
	sql := "INSERT INTO `piano`.`node_billing_info` (`group_id`, `date`, `node_uuid`, `default_charge_cpu`, `default_charge_memory`, `default_charge_nic`, `discount_rate`) VALUES "

	for _, info := range *infoList {
		sql += fmt.Sprintf("(%d, DATE(NOW()), %s, %f, %f, %f, %f),",
			info.GroupID,
			info.NodeUUID,
			info.DefChargeCPU,
			info.DefChargeMEM,
			info.DefChargeNIC,
			info.DiscountRate)
	}

	sql = strings.TrimSuffix(sql, ",") + " AS `new_info` "
	sql += "ON DUPLICATE KEY UPDATE " +
		"`default_charge_cpu` = `new_info`.`default_charge_cpu`, " +
		"`default_charge_memory` = `new_info`.`default_charge_memory`, " +
		"`default_charge_nic` = `new_info`.`default_charge_nic`, " +
		"`discount_rate` = `new_info`.`discount_rate`;"

	res, err := sendQuery(sql)
	res.Close()

	return err
}

func InsertServerBillingInfo(infoList *[]model.ServerBill) *errors.HccError {
	sql := "INSERT INTO `piano`.`server_billing_info` (`group_id`, `date`, `server_uuid`, `network_traffic`, `traffic_charge_per_kb`, `discount_rate`) VALUES "

	for _, info := range *infoList {
		sql += fmt.Sprintf("(%d, DATE(NOW()), %s, %d, %f, %f),",
			info.GroupID,
			info.ServerUUID,
			info.NetworkTraffic,
			info.TrafficChargePerKB,
			info.DiscountRate)
	}

	sql = strings.TrimSuffix(sql, ",") + " AS `new_info` "
	sql += "ON DUPLICATE KEY UPDATE " +
		"`network_traffic` = `new_info`.`network_traffic`, " +
		"`traffic_charge_per_kb` = `new_info`.`traffic_charge_per_kb`, " +
		"`discount_rate` = `new_info`.`discount_rate`;"

	res, err := sendQuery(sql)
	res.Close()

	return err
}

func InsertVolumeBillingInfo(infoList *[]model.VolumeBill) *errors.HccError {
	sql := "INSERT INTO `piano`.`volume_billing_info` (`group_id`, `date`, `hdd_size`, `ssd_size`, `nvme_size`, `hdd_charge_per_gb`, `ssd_charge_per_gb`, `nvme_charge_per_gb`, `discount_rate`) VALUES "

	for _, info := range *infoList {
		sql += fmt.Sprintf("(%d, DATE(NOW()), %d, %d, %d, %f, %f, %f, %f),",
			info.GroupID,
			info.HDDSize,
			info.SSDSize,
			info.NVMESize,
			info.HDDChargePerGB,
			info.SSDChargePerGB,
			info.NVMEChargePerGB,
			info.DiscountRate)
	}

	sql = strings.TrimSuffix(sql, ",") + " AS `new_info` "
	sql += "ON DUPLICATE KEY UPDATE " +
		"`hdd_size` = `new_info`.`hdd_size`," +
		"`ssd_size` = `new_info`.`ssd_size`," +
		"`nvme_size` = `new_info`.`nvme_size`," +
		"`hdd_charge_per_gb` = `new_info`.`hdd_charge_per_gb`," +
		"`ssd_charge_per_gb` = `new_info`.`ssd_charge_per_gb`," +
		"`nvme_charge_per_gb` = `new_info`.`nvme_charge_per_gb`," +
		"`discount_rate` = `new_info`.`discount_rate`;"

	res, err := sendQuery(sql)
	res.Close()

	return err
}

func InsertDailyInfo() *errors.HccError {

	sql := "INSERT INTO `piano`.`daily_info` (`date`, `group_id`, `charge_node`, `charge_server`, `charge_network`, `charge_volume`) SELECT	`current_info`.`date`, `current_info`.`group_id`, `current_info`.`charge_node`,	`current_info`.`charge_network`, `current_info`.`charge_server`, `current_info`.`charge_volume` FROM (SELECT `node_charge`.`date` AS `date`, `node_charge`.`group_id` AS `group_id`, `charge_node`, `charge_network`, `charge_server`, `charge_volume` FROM (SELECT `node`.`group_id`, `node`.`date`, SUM((`node`.`default_charge_cpu` + `node`.`default_charge_memory` + `default_charge_nic`) * (1 - `node`.`discount_rate`)) AS `charge_node` FROM `piano`.`node_billing_info` AS `node` WHERE `node`.`date` = DATE(NOW()) GROUP BY `node`.`group_id`, `node`.`date`) AS `node_charge` LEFT JOIN (SELECT `net`.`group_id`, `net`.`date`, (`net`.`subnet_count` * `net`.`subnet_charge_per_cnt` * (1-`net`.`discount_rate`) + `net`.`adaptiveip_count` * `net`.`adaptiveip_charge_per_cnt` * (1-`net`.`discount_rate`)) AS `charge_network` FROM `piano`.`network_billing_info` AS `net` WHERE `net`.`date` = DATE(NOW())) AS `net_charge` ON `net_charge`.`date` = `node_charge`.`date` AND `net_charge`.`group_id` = `node_charge`.`group_id` LEFT JOIN (SELECT `server`.`group_id`, `server`.`date`, SUM((`server`.`network_traffic` * `server`.`traffic_charge_per_kb`) * (1 - `server`.`discount_rate`)) AS `charge_server` FROM `piano`.`server_billing_info` AS `server` WHERE `server`.`date` = DATE(NOW()) GROUP BY `server`.`group_id`, `server`.`date`) AS `server_charge` ON `server_charge`.`date` = `node_charge`.`date` AND `server_charge`.`group_id` = `node_charge`.`group_id` LEFT JOIN (SELECT `volume`.`group_id`, `volume`.`date`, ((`volume`.`hdd_size` * `volume`.`hdd_charge_per_gb`) + (`volume`.`ssd_size` * `volume`.`ssd_charge_per_gb`) + (`volume`.`nvme_size` * `volume`.`nvme_charge_per_gb`)) * (1 - `volume`.`discount_rate`) AS `charge_volume` FROM `piano`.`volume_billing_info` AS `volume` WHERE `volume`.`date` = DATE(NOW())) AS `volume_charge` ON `volume_charge`.`date` = `node_charge`.`date` AND `volume_charge`.`group_id` = `node_charge`.`group_id`) AS `current_info` ON DUPLICATE KEY UPDATE `daily_info`.`charge_node` = `current_info`.`charge_node`, `daily_info`.`charge_server` = `current_info`.`charge_server`, `daily_info`.`charge_network` = `current_info`.`charge_network`, `daily_info`.`charge_volume` = `current_info`.`charge_volume`;"

	res, err := sendQuery(sql)
	res.Close()

	return err
}

func GetBill(groupID int, start, end, billType string) (*mysql.Rows, *errors.HccError) {
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
		sql = "SELECT * FROM `piano`.`" + billType + "_bill` WHERE `bill_id` BETWEEN " + billIdStart + " AND " + billIdEnd + ";"
	default:
		return nil, errors.NewHccError(errors.PianoSQLOperationFail, "DAO(GetBill) -> Unsupport billing type")
	}

	res, err := sendQuery(sql)

	return res, err
}

func GetBillInfo(groupID int, date, billType, category string) (*mysql.Rows, *errors.HccError) {
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
		return nil, errors.NewHccError(errors.PianoSQLOperationFail, "DAO(GetBillInfo) -> Unsupport billing type")
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
		return nil, errors.NewHccError(errors.PianoSQLOperationFail, "DAO(GetBillInfo) -> Unsupport category")
	}

	sql := "SELECT * FROM `piano`.`" + category + "_billing_info` WHERE `group_id`=" + strconv.Itoa(groupID) + " AND `date` BETWEEN DATE(" + dateStart + ") AND DATE(" + strconv.Itoa(dateEnd) + ");"

	res, err := sendQuery(sql)

	return res, err
}
