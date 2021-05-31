/* daily_bill */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`daily_bill` AS
SELECT
    CONCAT(LPAD(CAST(FLOOR((CAST(`piano`.`daily_info`.`date` AS UNSIGNED) / 10000))
                    AS CHAR (4) CHARSET UTF8MB4),
                4,
                0),
           '-',
           LPAD(CAST(FLOOR(((CAST(`piano`.`daily_info`.`date` AS UNSIGNED) % 10000) / 100))
                    AS CHAR (2) CHARSET UTF8MB4),
                2,
                0),
           '-',
           LPAD(CAST(FLOOR((CAST(`piano`.`daily_info`.`date` AS UNSIGNED) % 100))
                    AS CHAR (2) CHARSET UTF8MB4),
                2,
                0)) AS `date`,
    `piano`.`daily_info`.`group_id` AS `group_id`,
    `piccolo`.`group`.`name` AS `group_name`,
    `piano`.`daily_info`.`charge_node` AS `charge_node`,
    `piano`.`daily_info`.`charge_server` AS `charge_server`,
    `piano`.`daily_info`.`charge_network` AS `charge_network`,
    `piano`.`daily_info`.`charge_volume` AS `charge_volume`
FROM
    (`piano`.`daily_info`
        JOIN `piccolo`.`group`)
WHERE
    (`piano`.`daily_info`.`group_id` = `piccolo`.`group`.`id`)

/* monthly_bill */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`monthly_bill` AS
SELECT
    `daily`.`date` AS `date`,
    `daily`.`group_id` AS `group_id`,
    `daily`.`group_name` AS `group_name`,
    SUM(`daily`.`charge_node`) AS `charge_node`,
    SUM(`daily`.`charge_server`) AS `charge_server`,
    SUM(`daily`.`charge_network`) AS `charge_network`,
    SUM(`daily`.`charge_volume`) AS `charge_volume`
FROM
    (SELECT
         CONCAT(LPAD(CAST(FLOOR((CAST(REPLACE(`piano`.`daily_bill`.`date`, '-', '')
                                     AS UNSIGNED) / 10000))
                         AS CHAR (4) CHARSET UTF8MB4), 4, 0), '-', LPAD(CAST(FLOOR(((CAST(REPLACE(`piano`.`daily_bill`.`date`, '-', '')
                                                                                         AS UNSIGNED) % 10000) / 100))
                                                                            AS CHAR (2) CHARSET UTF8MB4), 2, 0)) AS `date`,
         `piano`.`daily_bill`.`group_id` AS `group_id`,
         `piano`.`daily_bill`.`group_name` AS `group_name`,
         `piano`.`daily_bill`.`charge_node` AS `charge_node`,
         `piano`.`daily_bill`.`charge_server` AS `charge_server`,
         `piano`.`daily_bill`.`charge_network` AS `charge_network`,
         `piano`.`daily_bill`.`charge_volume` AS `charge_volume`
     FROM
         `piano`.`daily_bill`) `daily`
GROUP BY `daily`.`date` , `daily`.`group_id` , `daily`.`group_name`

/* yearly_bill */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`yearly_bill` AS
SELECT
    `monthly`.`date` AS `date`,
    `monthly`.`group_id` AS `group_id`,
    `monthly`.`group_name` AS `group_name`,
    SUM(`monthly`.`charge_node`) AS `charge_node`,
    SUM(`monthly`.`charge_server`) AS `charge_server`,
    SUM(`monthly`.`charge_network`) AS `charge_network`,
    SUM(`monthly`.`charge_volume`) AS `charge_volume`
FROM
    (SELECT
         CONCAT(LPAD(CAST(FLOOR((CAST(REPLACE(`piano`.`monthly_bill`.`date`, '-', '')
                                     AS UNSIGNED) / 100))
                         AS CHAR (4) CHARSET UTF8MB4), 4, 0)) AS `date`,
         `piano`.`monthly_bill`.`group_id` AS `group_id`,
         `piano`.`monthly_bill`.`group_name` AS `group_name`,
         `piano`.`monthly_bill`.`charge_node` AS `charge_node`,
         `piano`.`monthly_bill`.`charge_server` AS `charge_server`,
         `piano`.`monthly_bill`.`charge_network` AS `charge_network`,
         `piano`.`monthly_bill`.`charge_volume` AS `charge_volume`
     FROM
         `piano`.`monthly_bill`) `monthly`
GROUP BY `monthly`.`date` , `monthly`.`group_id` , `monthly`.`group_name`

/* monthly_network_billing_info */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`monthly_network_billing_info` AS
SELECT
    `daily`.`group_id` AS `group_id`,
    `daily`.`date` AS `date`,
    SUM(`daily`.`charge_subnet`) AS `charge_subnet`,
    SUM(`daily`.`charge_adaptive_ip`) AS `charge_adaptive_ip`
FROM
    (SELECT
         `piano`.`network_billing_info`.`group_id` AS `group_id`,
         CONCAT(LPAD(CAST(FLOOR((CAST(REPLACE(`piano`.`network_billing_info`.`date`, '-', '')
                                     AS UNSIGNED) / 10000))
                         AS CHAR (4) CHARSET UTF8MB4), 4, 0), '-', LPAD(CAST(FLOOR(((CAST(REPLACE(`piano`.`network_billing_info`.`date`, '-', '')
                                                                                         AS UNSIGNED) % 10000) / 100))
                                                                            AS CHAR (2) CHARSET UTF8MB4), 2, 0)) AS `date`,
         `piano`.`network_billing_info`.`charge_subnet` AS `charge_subnet`,
         `piano`.`network_billing_info`.`charge_adaptive_ip` AS `charge_adaptive_ip`
     FROM
         `piano`.`network_billing_info`) `daily`
GROUP BY `daily`.`group_id` , `daily`.`date`

/* monthly_node_billing_info */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`monthly_node_billing_info` AS
SELECT
    `daily`.`group_id` AS `group_id`,
    `daily`.`date` AS `date`,
    `daily`.`node_uuid` AS `node_uuid`,
    SUM(`daily`.`charge_cpu`) AS `charge_cpu`,
    SUM(`daily`.`charge_memory`) AS `charge_memory`,
    SUM(`daily`.`charge_nic`) AS `charge_nic`
FROM
    (SELECT
         `piano`.`node_billing_info`.`group_id` AS `group_id`,
         CONCAT(LPAD(CAST(FLOOR((CAST(REPLACE(`piano`.`node_billing_info`.`date`, '-', '')
                                     AS UNSIGNED) / 10000))
                         AS CHAR (4) CHARSET UTF8MB4), 4, 0), '-', LPAD(CAST(FLOOR(((CAST(REPLACE(`piano`.`node_billing_info`.`date`, '-', '')
                                                                                         AS UNSIGNED) % 10000) / 100))
                                                                            AS CHAR (2) CHARSET UTF8MB4), 2, 0)) AS `date`,
         `piano`.`node_billing_info`.`node_uuid` AS `node_uuid`,
         `piano`.`node_billing_info`.`charge_cpu` AS `charge_cpu`,
         `piano`.`node_billing_info`.`charge_memory` AS `charge_memory`,
         `piano`.`node_billing_info`.`charge_nic` AS `charge_nic`
     FROM
         `piano`.`node_billing_info`) `daily`
GROUP BY `daily`.`group_id` , `daily`.`date` , `daily`.`node_uuid`

/* monthly_server_billing_info */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`monthly_server_billing_info` AS
SELECT
    `daily`.`group_id` AS `group_id`,
    `daily`.`date` AS `date`,
    `daily`.`server_uuid` AS `server_uuid`,
    SUM(`daily`.`charge_traffic`) AS `charge_traffic`
FROM
    (SELECT
         `piano`.`server_billing_info`.`group_id` AS `group_id`,
         CONCAT(LPAD(CAST(FLOOR((CAST(REPLACE(`piano`.`server_billing_info`.`date`, '-', '')
                                     AS UNSIGNED) / 10000))
                         AS CHAR (4) CHARSET UTF8MB4), 4, 0), '-', LPAD(CAST(FLOOR(((CAST(REPLACE(`piano`.`server_billing_info`.`date`, '-', '')
                                                                                         AS UNSIGNED) % 10000) / 100))
                                                                            AS CHAR (2) CHARSET UTF8MB4), 2, 0)) AS `date`,
         `piano`.`server_billing_info`.`server_uuid` AS `server_uuid`,
         `piano`.`server_billing_info`.`charge_traffic` AS `charge_traffic`
     FROM
         `piano`.`server_billing_info`) `daily`
GROUP BY `daily`.`group_id` , `daily`.`date` , `daily`.`server_uuid`

/* monthly_volume_billing_info */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`monthly_volume_billing_info` AS
SELECT
    `daily`.`group_id` AS `group_id`,
    `daily`.`date` AS `date`,
    SUM(`daily`.`charge_ssd`) AS `charge_ssd`,
    SUM(`daily`.`charge_hdd`) AS `charge_hdd`
FROM
    (SELECT
         `piano`.`volume_billing_info`.`group_id` AS `group_id`,
         CONCAT(LPAD(CAST(FLOOR((CAST(REPLACE(`piano`.`volume_billing_info`.`date`, '-', '')
                                     AS UNSIGNED) / 10000))
                         AS CHAR (4) CHARSET UTF8MB4), 4, 0), '-', LPAD(CAST(FLOOR(((CAST(REPLACE(`piano`.`volume_billing_info`.`date`, '-', '')
                                                                                         AS UNSIGNED) % 10000) / 100))
                                                                            AS CHAR (2) CHARSET UTF8MB4), 2, 0)) AS `date`,
         `piano`.`volume_billing_info`.`charge_ssd` AS `charge_ssd`,
         `piano`.`volume_billing_info`.`charge_hdd` AS `charge_hdd`
     FROM
         `piano`.`volume_billing_info`) `daily`
GROUP BY `daily`.`group_id` , `daily`.`date`

/* yearly_network_billing_info */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`yearly_network_billing_info` AS
SELECT
    `monthly`.`group_id` AS `group_id`,
    `monthly`.`date` AS `date`,
    SUM(`monthly`.`charge_subnet`) AS `charge_subnet`,
    SUM(`monthly`.`charge_adaptive_ip`) AS `charge_adaptive_ip`
FROM
    (SELECT
         `piano`.`monthly_network_billing_info`.`group_id` AS `group_id`,
         CONCAT(LPAD(CAST(FLOOR((CAST(REPLACE(`piano`.`monthly_network_billing_info`.`date`, '-', '')
                                     AS UNSIGNED) / 100))
                         AS CHAR (4) CHARSET UTF8MB4), 4, 0)) AS `date`,
         `piano`.`monthly_network_billing_info`.`charge_subnet` AS `charge_subnet`,
         `piano`.`monthly_network_billing_info`.`charge_adaptive_ip` AS `charge_adaptive_ip`
     FROM
         `piano`.`monthly_network_billing_info`) `monthly`
GROUP BY `monthly`.`group_id` , `monthly`.`date`

/* yearly_node_billing_info */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`yearly_node_billing_info` AS
SELECT
    `monthly`.`group_id` AS `group_id`,
    `monthly`.`date` AS `date`,
    `monthly`.`node_uuid` AS `node_uuid`,
    SUM(`monthly`.`charge_cpu`) AS `charge_cpu`,
    SUM(`monthly`.`charge_memory`) AS `charge_memory`,
    SUM(`monthly`.`charge_nic`) AS `charge_nic`
FROM
    (SELECT
         `piano`.`monthly_node_billing_info`.`group_id` AS `group_id`,
         CONCAT(LPAD(CAST(FLOOR((CAST(REPLACE(`piano`.`monthly_node_billing_info`.`date`, '-', '')
                                     AS UNSIGNED) / 100))
                         AS CHAR (4) CHARSET UTF8MB4), 4, 0)) AS `date`,
         `piano`.`monthly_node_billing_info`.`node_uuid` AS `node_uuid`,
         `piano`.`monthly_node_billing_info`.`charge_cpu` AS `charge_cpu`,
         `piano`.`monthly_node_billing_info`.`charge_memory` AS `charge_memory`,
         `piano`.`monthly_node_billing_info`.`charge_nic` AS `charge_nic`
     FROM
         `piano`.`monthly_node_billing_info`) `monthly`
GROUP BY `monthly`.`group_id` , `monthly`.`date` , `monthly`.`node_uuid`

/* yearly_server_billing_info */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`yearly_server_billing_info` AS
SELECT
    `monthly`.`group_id` AS `group_id`,
    `monthly`.`date` AS `date`,
    `monthly`.`server_uuid` AS `server_uuid`,
    SUM(`monthly`.`charge_traffic`) AS `charge_traffic`
FROM
    (SELECT
         `piano`.`monthly_server_billing_info`.`group_id` AS `group_id`,
         CONCAT(LPAD(CAST(FLOOR((CAST(REPLACE(`piano`.`monthly_server_billing_info`.`date`, '-', '')
                                     AS UNSIGNED) / 100))
                         AS CHAR (4) CHARSET UTF8MB4), 4, 0)) AS `date`,
         `piano`.`monthly_server_billing_info`.`server_uuid` AS `server_uuid`,
         `piano`.`monthly_server_billing_info`.`charge_traffic` AS `charge_traffic`
     FROM
         `piano`.`monthly_server_billing_info`) `monthly`
GROUP BY `monthly`.`group_id` , `monthly`.`date` , `monthly`.`server_uuid`

/* yearly_volume_billing_info */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`yearly_volume_billing_info` AS
SELECT
    `monthly`.`group_id` AS `group_id`,
    `monthly`.`date` AS `date`,
    SUM(`monthly`.`charge_ssd`) AS `charge_ssd`,
    SUM(`monthly`.`charge_hdd`) AS `charge_hdd`
FROM
    (SELECT
         `piano`.`monthly_volume_billing_info`.`group_id` AS `group_id`,
         CONCAT(LPAD(CAST(FLOOR((CAST(REPLACE(`piano`.`monthly_volume_billing_info`.`date`, '-', '')
                                     AS UNSIGNED) / 100))
                         AS CHAR (4) CHARSET UTF8MB4), 4, 0)) AS `date`,
         `piano`.`monthly_volume_billing_info`.`charge_ssd` AS `charge_ssd`,
         `piano`.`monthly_volume_billing_info`.`charge_hdd` AS `charge_hdd`
     FROM
         `piano`.`monthly_volume_billing_info`) `monthly`
GROUP BY `monthly`.`group_id` , `monthly`.`date`
