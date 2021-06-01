package dao

import (
	"errors"
	"fmt"
	"hcc/piano/lib/config"
	"hcc/piano/lib/logger"
	"innogrid.com/hcloud-classic/pb"
	"strconv"
	"strings"
	"time"

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

func GetDailyInfo(groupList []*pb.Group,
	nodeBillingList *[]model.NodeBill,
	serverBillingList *[]model.ServerBill,
	networkBillingList *[]model.NetworkBill,
	volumeBillingList *[]model.VolumeBill) *[]model.DailyBill {
	var billList []model.DailyBill

	if nodeBillingList == nil {
		nodeBillingList = &[]model.NodeBill{}
	}
	if serverBillingList == nil {
		serverBillingList = &[]model.ServerBill{}
	}
	if networkBillingList == nil {
		networkBillingList = &[]model.NetworkBill{}
	}
	if volumeBillingList == nil {
		volumeBillingList = &[]model.VolumeBill{}
	}

	for _, group := range groupList {
		if group.Id == 1 {
			continue
		}

		var chargeNode int64 = 0
		var chargeServer int64 = 0
		var chargeNetwork int64 = 0
		var chargeVolume int64 = 0

		for _, nodeBilling := range *nodeBillingList {
			if nodeBilling.GroupID == group.Id {
				chargeNode += nodeBilling.ChargeCPU +
					nodeBilling.ChargeMEM +
					nodeBilling.ChargeNIC
			}
		}

		for _, serverBilling := range *serverBillingList {
			if serverBilling.GroupID == group.Id {
				chargeServer += serverBilling.ChargeTraffic
			}
		}

		for _, networkBilling := range *networkBillingList {
			if networkBilling.GroupID == group.Id {
				chargeNetwork += networkBilling.ChargeSubnet +
					networkBilling.ChargeAdaptiveIP
			}
		}

		for _, volumeBilling := range *volumeBillingList {
			if volumeBilling.GroupID == group.Id {
				chargeVolume += volumeBilling.ChargeSSD +
					volumeBilling.ChargeHDD
			}
		}

		billList = append(billList, model.DailyBill{
			GroupID:       group.Id,
			ChargeNode:    chargeNode,
			ChargeServer:  chargeServer,
			ChargeNetwork: chargeNetwork,
			ChargeVolume:  chargeVolume,
		})
	}

	return &billList
}

func InsertDailyInfo(infoList *[]model.DailyBill) error {
	sql := "INSERT INTO `piano`.`daily_info` (`group_id`, `date`, `charge_node`, `charge_server`, `charge_network`, `charge_volume`) VALUES "

	for _, info := range *infoList {
		sql += fmt.Sprintf("(%d, DATE(NOW()), %d, %d, %d, %d),",
			info.GroupID,
			info.ChargeNode,
			info.ChargeServer,
			info.ChargeNetwork,
			info.ChargeVolume)
	}

	sql = strings.TrimSuffix(sql, ",") + " AS `new_info` "
	sql += "ON DUPLICATE KEY UPDATE " +
		"`charge_node` = `new_info`.`charge_node`, " +
		"`charge_server` = `new_info`.`charge_server`, " +
		"`charge_network` = `new_info`.`charge_network`, " +
		"`charge_volume` = `new_info`.`charge_volume`;"

	res, err := sendQuery(sql)
	if res != nil {
		_ = res.Close()
	}

	return err
}

func GetBill(groupID *[]int64, start, end, billType string, row, page int) (*dbsql.Rows, error) {
	var dateStart string
	var dateEnd string
	billType = strings.ToLower(billType)
	sql := ""
	var groupIDQuery string

	if row == 0 || page == 0 {
		return nil, errors.New("need row and page arguments")
	}

	currentTime := time.Now()
	yyFront := currentTime.Format("2006")[:2]

	switch billType {
	case "daily":
		dateStart = yyFront + start[:2] + "-" + start[2:4] + "-" + start[4:6]
		dateEnd = yyFront + end[:2] + "-" + end[2:4] + "-" + end[4:6]
	case "monthly":
		dateStart = yyFront + start[:2] + "-" + start[2:4]
		dateEnd = yyFront + end[:2] + "-" + end[2:4]
	case "yearly":
		dateStart = yyFront + start[:2]
		dateEnd = yyFront + end[:2]
	default:
		return nil, errors.New("DAO(GetBill) -> Unsupported billing type")
	}

	if len(*groupID) != 0 {
		groupIDQuery = "AND ("
		for i := range *groupID {
			if i == 0 {
				groupIDQuery += "`group_id` = " + strconv.Itoa(int((*groupID)[i]))
				continue
			}
			groupIDQuery += " OR `group_id` = " + strconv.Itoa(int((*groupID)[i]))
		}
		groupIDQuery += ")"
	} else {
		groupIDQuery = ""
	}

	sql = "SELECT * FROM `piano`.`" + billType + "_bill` WHERE `date` BETWEEN '" + dateStart + "' AND '" +
		dateEnd + "' " + groupIDQuery + " ORDER BY `date` DESC, `group_id` ASC " +
		"LIMIT " + strconv.Itoa(row) + " OFFSET " + strconv.Itoa(row*(page-1)) + ";"

	if config.Billing.Debug == "on" {
		logger.Logger.Println("Sending SQL Query from GetBill(): " + sql)
	}
	res, err := sendQuery(sql)

	return res, err
}

func GetBillCount(groupID *[]int64, start, end, billType string) (int, error) {
	var billCount int

	var dateStart string
	var dateEnd string
	billType = strings.ToLower(billType)
	sql := ""
	var groupIDQuery string

	currentTime := time.Now()
	yyFront := currentTime.Format("2006")[:2]

	switch billType {
	case "daily":
		dateStart = yyFront + start[:2] + "-" + start[2:4] + "-" + start[4:6]
		dateEnd = yyFront + end[:2] + "-" + end[2:4] + "-" + end[4:6]
	case "monthly":
		dateStart = yyFront + start[:2] + "-" + start[2:4]
		dateEnd = yyFront + end[:2] + "-" + end[2:4]
	case "yearly":
		dateStart = yyFront + start[:2]
		dateEnd = yyFront + end[:2]
	default:
		return 0, errors.New("DAO(GetBill) -> Unsupported billing type")
	}

	if len(*groupID) != 0 {
		groupIDQuery = "AND ("
		for i := range *groupID {
			if i == 0 {
				groupIDQuery += "`group_id` = " + strconv.Itoa(int((*groupID)[i]))
				continue
			}
			groupIDQuery += " OR `group_id` = " + strconv.Itoa(int((*groupID)[i]))
		}
		groupIDQuery += ")"
	} else {
		groupIDQuery = ""
	}

	sql = "SELECT COUNT(*) FROM `piano`.`" + billType + "_bill` WHERE `date` BETWEEN '" + dateStart + "' AND '" +
		dateEnd + "' " + groupIDQuery + ";"

	if config.Billing.Debug == "on" {
		logger.Logger.Println("Sending SQL Query from GetBillCount(): " + sql)
	}
	res, err := sendQuery(sql)
	if err != nil {
		return 0, err
	}

	res.Next()
	err = res.Scan(&billCount)
	if err != nil {
		return 0, err
	}

	return billCount, nil
}

func GetBillInfo(groupID int64, date, billType, category string) (*dbsql.Rows, error) {
	billType = strings.ToLower(billType)
	category = strings.ToLower(category)
	var prefix string

	currentTime := time.Now()
	yyFront := currentTime.Format("2006")[:2]

	switch billType {
	case "daily":
		date = yyFront + date[:2] + "-" + date[2:4] + "-" + date[4:6]
	case "monthly":
		date = yyFront + date[:2] + "-" + date[2:4]
	case "yearly":
		date = yyFront + date[:2]
	default:
		return nil, errors.New("DAO(GetBill) -> Unsupported billing type")
	}

	prefix = billType + "_" + category

	switch category {
	case "node":
		fallthrough
	case "server":
		break
	case "network":
		fallthrough
	case "volume":
		if billType == "daily" {
			prefix = category
		}
		break
	default:
		return nil, errors.New("DAO(GetBill) -> Unsupported category")
	}

	sql := "SELECT * FROM `piano`.`" + prefix + "_billing_info` WHERE `group_id`=" + strconv.Itoa(int(groupID)) + " AND `date`='" + date + "';"

	if config.Billing.Debug == "on" {
		logger.Logger.Println("Sending SQL Query from GetBillInfo(): " + sql)
	}
	res, err := sendQuery(sql)

	return res, err
}
