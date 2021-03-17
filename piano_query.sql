
# NET
SELECT * FROM `piano`.`network_billing_info`;
DELETE FROM `piano`.`network_billing_info` WHERE `group_id` > 0;
INSERT INTO `piano`.`network_billing_info`
(`group_id`, `date`, `subnet_count`, `adaptiveip_count`, `subnet_charge_per_cnt`, `adaptiveip_charge_per_cnt`, `discount_rate`)
VALUES
(1000, DATE('2021-04-01'), 2, 2, 1000, 10000, 0),
(1001, DATE('2021-04-01'), 1, 1, 1000, 10000, 0),
(1002, DATE('2021-04-01'), 1, 1, 1000, 10000, 0),
(1003, DATE('2021-04-01'), 1, 1, 1000, 10000, 0) AS `new_info`
ON DUPLICATE KEY
UPDATE
	`subnet_count` = `new_info`.`subnet_count`,
    `adaptiveip_count` = `new_info`.`adaptiveip_count`,
    `subnet_charge_per_cnt` = `new_info`.`subnet_charge_per_cnt`,
    `adaptiveip_charge_per_cnt` = `new_info`.`adaptiveip_charge_per_cnt`,
    `discount_rate` = `new_info`.`discount_rate`;

# NODE
SELECT * FROM `piano`.`node_billing_info`;
DELETE FROM `piano`.`node_billing_info` WHERE `group_id` > 0;
INSERT INTO `piano`.`node_billing_info`
(`group_id`, `date`, `node_uuid`, `default_charge_cpu`, `default_charge_memory`, `default_charge_nic`, `discount_rate`)
VALUES
(1000, DATE('2021-04-01'), "1000001", 100, 10, 1, 0),
(1000, DATE('2021-04-01'), "1000002", 100, 10, 1, 0),
(1001, DATE('2021-04-01'), "1001001", 100, 10, 1, 0),
(1002, DATE('2021-04-01'), "1002001", 100, 10, 1, 0),
(1003, DATE('2021-04-01'), "1003001", 100, 10, 1, 0) AS `new_info`
ON DUPLICATE KEY
UPDATE
	`default_charge_cpu` = `new_info`.`default_charge_cpu`,
    `default_charge_memory` = `new_info`.`default_charge_memory`,
    `default_charge_nic` = `new_info`.`default_charge_nic`,
    `discount_rate` = `new_info`.`discount_rate`;


# SERVER
SELECT * FROM `piano`.`server_billing_info`;
DELETE FROM `piano`.`server_billing_info` WHERE `group_id` > 0;
INSERT INTO `piano`.`server_billing_info`
(`group_id`, `date`, `server_uuid`, `network_traffic`, `traffic_charge_per_kb`, `discount_rate`)
VALUES
(1000, DATE('2021-04-01'), "1000100", 100000, 1, 0),
(1000, DATE('2021-04-01'), "1000200", 100000, 1, 0),
(1001, DATE('2021-04-01'), "1001100", 100000, 1, 0),
(1002, DATE('2021-04-01'), "1002100", 100000, 1, 0),
(1003, DATE('2021-04-01'), "1003100", 100000, 1, 0) AS `new_info`
ON DUPLICATE KEY
UPDATE
	`network_traffic` = `new_info`.`network_traffic`,
    `traffic_charge_per_kb` = `new_info`.`traffic_charge_per_kb`,
    `discount_rate` = `new_info`.`discount_rate`;


# VOLUME
SELECT * FROM `piano`.`volume_billing_info`;
DELETE FROM `piano`.`volume_billing_info` WHERE `group_id` > 0;
INSERT INTO `piano`.`volume_billing_info`
(`group_id`, `date`, `hdd_size`, `ssd_size`, `nvme_size`, `hdd_charge_per_gb`, `ssd_charge_per_gb`, `nvme_charge_per_gb`, `discount_rate`)
VALUES
(1000, DATE('2021-04-01'), 100, 10, 1, 1000, 1000, 1000, 0),
(1001, DATE('2021-04-01'), 100, 10, 1, 1000, 1000, 1000, 0),
(1002, DATE('2021-04-01'), 100, 10, 1, 1000, 1000, 1000, 0),
(1003, DATE('2021-04-01'), 100, 10, 1, 1000, 1000, 1000, 0) AS `new_info`
ON DUPLICATE KEY
UPDATE
	`hdd_size` = `new_info`.`hdd_size`,
    `ssd_size` = `new_info`.`ssd_size`,
	`nvme_size` = `new_info`.`nvme_size`,
    `hdd_charge_per_gb` = `new_info`.`hdd_charge_per_gb`,
    `ssd_charge_per_gb` = `new_info`.`ssd_charge_per_gb`,
    `nvme_charge_per_gb` = `new_info`.`nvme_charge_per_gb`,
    `discount_rate` = `new_info`.`discount_rate`;

SELECT * FROM `piano`.`daily_info`;
SELECT * FROM `piano`.`daily_bill`; -- View
SELECT * FROM `piano`.`monthly_bill`; -- View
SELECT * FROM `piano`.`yearly_bill`; -- View
DELETE FROM `piano`.`daily_info` WHERE `bill_id` > 0;

INSERT INTO `piano`.`daily_info` (`date`, `group_id`, `charge_node`, `charge_server`, `charge_network`, `charge_volume`)
SELECT
	`current_info`.`date`,
    `current_info`.`group_id`,
    `current_info`.`charge_node`,
    `current_info`.`charge_network`,
    `current_info`.`charge_server`,
    `current_info`.`charge_volume`
FROM
	(SELECT
		`node_charge`.`date` AS `date`,
		`node_charge`.`group_id` AS `group_id`,
		`charge_node`,
		`charge_network`,
		`charge_server`,
		`charge_volume`
	FROM
		(SELECT
			`node`.`group_id`,
			`node`.`date`,
			SUM((`node`.`default_charge_cpu` + `node`.`default_charge_memory` + `default_charge_nic`) * (1 - `node`.`discount_rate`)) AS `charge_node`
		FROM
			`piano`.`node_billing_info` AS `node`
		WHERE
			`node`.`date` BETWEEN DATE('2021-03-10') AND DATE('2021-04-01')
		GROUP BY
			`node`.`group_id`, `node`.`date`) AS `node_charge`
		LEFT JOIN
			(SELECT
				`net`.`group_id`,
				`net`.`date`,
				(`net`.`subnet_count` * `net`.`subnet_charge_per_cnt` * (1-`net`.`discount_rate`) +
				 `net`.`adaptiveip_count` * `net`.`adaptiveip_charge_per_cnt` * (1-`net`.`discount_rate`)) AS `charge_network`
			FROM
				`piano`.`network_billing_info` AS `net`
			WHERE
				`net`.`date` BETWEEN DATE('2021-03-10') AND DATE('2021-04-01')) AS `net_charge`
		ON `net_charge`.`date` = `node_charge`.`date` AND `net_charge`.`group_id` = `node_charge`.`group_id`
		LEFT JOIN
			(SELECT
				`server`.`group_id`,
				`server`.`date`,
				SUM((`server`.`network_traffic` * `server`.`traffic_charge_per_kb`) * (1 - `server`.`discount_rate`)) AS `charge_server`
			FROM
				`piano`.`server_billing_info` AS `server`
			WHERE
				`server`.`date` BETWEEN DATE('2021-03-10') AND DATE('2021-04-01')
			GROUP BY
				`server`.`group_id`, `server`.`date`) AS `server_charge`
		ON `server_charge`.`date` = `node_charge`.`date` AND `server_charge`.`group_id` = `node_charge`.`group_id`
		LEFT JOIN
			(SELECT
				`volume`.`group_id`,
				`volume`.`date`,
				((`volume`.`hdd_size` * `volume`.`hdd_charge_per_gb`) +
				 (`volume`.`ssd_size` * `volume`.`ssd_charge_per_gb`) +
				 (`volume`.`nvme_size` * `volume`.`nvme_charge_per_gb`)) * (1 - `volume`.`discount_rate`) AS `charge_volume`
			FROM
				`piano`.`volume_billing_info` AS `volume`
			WHERE
				`volume`.`date` BETWEEN DATE('2021-03-10') AND DATE('2021-04-01')) AS `volume_charge`
		ON `volume_charge`.`date` = `node_charge`.`date` AND `volume_charge`.`group_id` = `node_charge`.`group_id`
	) AS `current_info`
ON DUPLICATE KEY UPDATE
	`daily_info`.`charge_node` = `current_info`.`charge_node`,
    `daily_info`.`charge_server` = `current_info`.`charge_server`,
    `daily_info`.`charge_network` = `current_info`.`charge_network`,
    `daily_info`.`charge_volume` = `current_info`.`charge_volume`;
