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
    (`piano`.`daily_info`.`group_id` = `piccolo`.`group`.`id`);

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
GROUP BY `daily`.`date` , `daily`.`group_id` , `daily`.`group_name`;

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
GROUP BY `monthly`.`date` , `monthly`.`group_id` , `monthly`.`group_name`;

/* daily_node_billing_info */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`daily_node_billing_info` AS
SELECT
    `daily`.`group_id` AS `group_id`,
    `daily`.`date` AS `date`,
    `daily`.`node_uuid` AS `node_uuid`,
    `daily`.`charge_cpu` AS `charge_cpu`,
    `daily`.`charge_memory` AS `charge_memory`,
    `daily`.`charge_nic` AS `charge_nic`,
    `daily`.`node_uptime_uptime_ms` AS `uptime_ms`
FROM
    (SELECT
         `piano`.`node_billing_info`.`group_id` AS `group_id`,
         `piano`.`node_billing_info`.`date` AS `date`,
         `piano`.`node_billing_info`.`node_uuid` AS `node_uuid`,
         `piano`.`node_billing_info`.`charge_cpu` AS `charge_cpu`,
         `piano`.`node_billing_info`.`charge_memory` AS `charge_memory`,
         `piano`.`node_billing_info`.`charge_nic` AS `charge_nic`,
         `flute`.`node_uptime`.`node_uuid` AS `node_uptime_node_uuid`,
         `flute`.`node_uptime`.`uptime_ms` AS `node_uptime_uptime_ms`,
         CAST(`flute`.`node_uptime`.`day` AS DATE) AS `node_uptime_date`
     FROM
         (`piano`.`node_billing_info`
             JOIN `flute`.`node_uptime`)) `daily`
WHERE
    ((`daily`.`node_uuid` = CONVERT( `daily`.`node_uptime_node_uuid` USING UTF8MB4))
        AND (`daily`.`date` = `daily`.`node_uptime_date`));

/* daily_server_billing_info */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`daily_server_billing_info` AS
SELECT
    `daily`.`group_id` AS `group_id`,
    `daily`.`date` AS `date`,
    `daily`.`server_uuid` AS `server_uuid`,
    `daily`.`charge_traffic` AS `charge_traffic`,
    `daily`.`traffic_kb` AS `traffic_kb`
FROM
    (SELECT
         `piano`.`server_billing_info`.`group_id` AS `group_id`,
         `piano`.`server_billing_info`.`date` AS `date`,
         `piano`.`server_billing_info`.`server_uuid` AS `server_uuid`,
         `piano`.`server_billing_info`.`charge_traffic` AS `charge_traffic`,
         CAST(`traffic`.`day` AS DATE) AS `traffic_date`,
         `traffic`.`server_uuid` AS `traffic_server_uuid`,
         (`traffic`.`rx_kb` + `traffic`.`tx_kb`) AS `traffic_kb`
     FROM
         (`piano`.`server_billing_info`
             JOIN `traffic`)) `daily`
WHERE
    ((`daily`.`traffic_date` = `daily`.`date`)
        AND (`daily`.`server_uuid` = CONVERT( `daily`.`traffic_server_uuid` USING UTF8MB4)));

/* monthly_subnet_billing_info */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`monthly_subnet_billing_info` AS
SELECT
    `daily`.`group_id` AS `group_id`,
    `daily`.`date` AS `date`,
    `daily`.`subnet_uuid` AS `subnet_uuid`,
    SUM(`daily`.`charge_subnet`) AS `charge_subnet`
FROM
    (SELECT
         `piano`.`subnet_billing_info`.`group_id` AS `group_id`,
         CONCAT(LPAD(CAST(FLOOR((CAST(REPLACE(`piano`.`subnet_billing_info`.`date`, '-', '')
                                     AS UNSIGNED) / 10000))
                         AS CHAR (4) CHARSET UTF8MB4), 4, 0), '-', LPAD(CAST(FLOOR(((CAST(REPLACE(`piano`.`subnet_billing_info`.`date`, '-', '')
                                                                                         AS UNSIGNED) % 10000) / 100))
                                                                            AS CHAR (2) CHARSET UTF8MB4), 2, 0)) AS `date`,
         `piano`.`subnet_billing_info`.`subnet_uuid` AS `subnet_uuid`,
         `piano`.`subnet_billing_info`.`charge_subnet` AS `charge_subnet`
     FROM
         `piano`.`subnet_billing_info`) `daily`
GROUP BY `daily`.`group_id` , `daily`.`date` , `daily`.`subnet_uuid`;

/* monthly_adaptiveip_billing_info */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`monthly_adaptiveip_billing_info` AS
SELECT
    `daily`.`group_id` AS `group_id`,
    `daily`.`date` AS `date`,
    `daily`.`server_uuid` AS `server_uuid`,
    SUM(`daily`.`charge_adaptiveip`) AS `charge_adaptiveip`
FROM
    (SELECT
         `piano`.`adaptiveip_billing_info`.`group_id` AS `group_id`,
         CONCAT(LPAD(CAST(FLOOR((CAST(REPLACE(`piano`.`adaptiveip_billing_info`.`date`, '-', '')
                                     AS UNSIGNED) / 10000))
                         AS CHAR (4) CHARSET UTF8MB4), 4, 0), '-', LPAD(CAST(FLOOR(((CAST(REPLACE(`piano`.`adaptiveip_billing_info`.`date`, '-', '')
                                                                                         AS UNSIGNED) % 10000) / 100))
                                                                            AS CHAR (2) CHARSET UTF8MB4), 2, 0)) AS `date`,
         `piano`.`adaptiveip_billing_info`.`server_uuid` AS `server_uuid`,
         `piano`.`adaptiveip_billing_info`.`charge_adaptiveip` AS `charge_adaptiveip`
     FROM
         `piano`.`adaptiveip_billing_info`) `daily`
GROUP BY `daily`.`group_id` , `daily`.`date` , `daily`.`server_uuid`;

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
    SUM(`daily`.`charge_nic`) AS `charge_nic`,
    SUM(`daily`.`uptime_ms`) AS `uptime_ms`
FROM
    (SELECT
         `piano`.`daily_node_billing_info`.`group_id` AS `group_id`,
         CONCAT(LPAD(CAST(FLOOR((CAST(REPLACE(`piano`.`daily_node_billing_info`.`date`, '-', '')
                                     AS UNSIGNED) / 10000))
                         AS CHAR (4) CHARSET UTF8MB4), 4, 0), '-', LPAD(CAST(FLOOR(((CAST(REPLACE(`piano`.`daily_node_billing_info`.`date`, '-', '')
                                                                                         AS UNSIGNED) % 10000) / 100))
                                                                            AS CHAR (2) CHARSET UTF8MB4), 2, 0)) AS `date`,
         `piano`.`daily_node_billing_info`.`node_uuid` AS `node_uuid`,
         `piano`.`daily_node_billing_info`.`charge_cpu` AS `charge_cpu`,
         `piano`.`daily_node_billing_info`.`charge_memory` AS `charge_memory`,
         `piano`.`daily_node_billing_info`.`charge_nic` AS `charge_nic`,
         `piano`.`daily_node_billing_info`.`uptime_ms` AS `uptime_ms`
     FROM
         `piano`.`daily_node_billing_info`) `daily`
GROUP BY `daily`.`group_id` , `daily`.`date` , `daily`.`node_uuid`;

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
    SUM(`daily`.`charge_traffic`) AS `charge_traffic`,
    SUM(`daily`.`traffic_kb`) AS `traffic_kb`
FROM
    (SELECT
         `piano`.`daily_server_billing_info`.`group_id` AS `group_id`,
         CONCAT(LPAD(CAST(FLOOR((CAST(REPLACE(`piano`.`daily_server_billing_info`.`date`, '-', '')
                                     AS UNSIGNED) / 10000))
                         AS CHAR (4) CHARSET UTF8MB4), 4, 0), '-', LPAD(CAST(FLOOR(((CAST(REPLACE(`piano`.`daily_server_billing_info`.`date`, '-', '')
                                                                                         AS UNSIGNED) % 10000) / 100))
                                                                            AS CHAR (2) CHARSET UTF8MB4), 2, 0)) AS `date`,
         `piano`.`daily_server_billing_info`.`server_uuid` AS `server_uuid`,
         `piano`.`daily_server_billing_info`.`charge_traffic` AS `charge_traffic`,
         `piano`.`daily_server_billing_info`.`traffic_kb` AS `traffic_kb`
     FROM
         `piano`.`daily_server_billing_info`) `daily`
GROUP BY `daily`.`group_id` , `daily`.`date` , `daily`.`server_uuid`;

/* monthly_volume_billing_info */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`monthly_volume_billing_info` AS
SELECT
    `daily`.`group_id` AS `group_id`,
    `daily`.`date` AS `date`,
    `daily`.`volume_uuid` AS `volume_uuid`,
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
         `piano`.`volume_billing_info`.`volume_uuid` AS `volume_uuid`,
         `piano`.`volume_billing_info`.`charge_ssd` AS `charge_ssd`,
         `piano`.`volume_billing_info`.`charge_hdd` AS `charge_hdd`
     FROM
         `piano`.`volume_billing_info`) `daily`
GROUP BY `daily`.`group_id` , `daily`.`date` , `daily`.`volume_uuid`;

/* yearly_subnet_billing_info */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`yearly_subnet_billing_info` AS
SELECT
    `monthly`.`group_id` AS `group_id`,
    `monthly`.`date` AS `date`,
    `monthly`.`subnet_uuid` AS `subnet_uuid`,
    SUM(`monthly`.`charge_subnet`) AS `charge_subnet`
FROM
    (SELECT
         `piano`.`monthly_subnet_billing_info`.`group_id` AS `group_id`,
         CONCAT(LPAD(CAST(FLOOR((CAST(REPLACE(`piano`.`monthly_subnet_billing_info`.`date`, '-', '')
                                     AS UNSIGNED) / 100))
                         AS CHAR (4) CHARSET UTF8MB4), 4, 0)) AS `date`,
         `piano`.`monthly_subnet_billing_info`.`subnet_uuid` AS `subnet_uuid`,
         `piano`.`monthly_subnet_billing_info`.`charge_subnet` AS `charge_subnet`
     FROM
         `piano`.`monthly_subnet_billing_info`) `monthly`
GROUP BY `monthly`.`group_id` , `monthly`.`date` , `monthly`.`subnet_uuid`;

/* yearly_adaptiveip_billing_info */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`yearly_adaptiveip_billing_info` AS
SELECT
    `monthly`.`group_id` AS `group_id`,
    `monthly`.`date` AS `date`,
    `monthly`.`server_uuid` AS `server_uuid`,
    SUM(`monthly`.`charge_adaptiveip`) AS `charge_adaptiveip`
FROM
    (SELECT
         `piano`.`monthly_adaptiveip_billing_info`.`group_id` AS `group_id`,
         CONCAT(LPAD(CAST(FLOOR((CAST(REPLACE(`piano`.`monthly_adaptiveip_billing_info`.`date`, '-', '')
                                     AS UNSIGNED) / 100))
                         AS CHAR (4) CHARSET UTF8MB4), 4, 0)) AS `date`,
         `piano`.`monthly_adaptiveip_billing_info`.`server_uuid` AS `server_uuid`,
         `piano`.`monthly_adaptiveip_billing_info`.`charge_adaptiveip` AS `charge_adaptiveip`
     FROM
         `piano`.`monthly_adaptiveip_billing_info`) `monthly`
GROUP BY `monthly`.`group_id` , `monthly`.`date` , `monthly`.`server_uuid`;

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
    SUM(`monthly`.`charge_nic`) AS `charge_nic`,
    SUM(`monthly`.`uptime_ms`) AS `uptime_ms`
FROM
    (SELECT
         `piano`.`monthly_node_billing_info`.`group_id` AS `group_id`,
         CONCAT(LPAD(CAST(FLOOR((CAST(REPLACE(`piano`.`monthly_node_billing_info`.`date`, '-', '')
                                     AS UNSIGNED) / 100))
                         AS CHAR (4) CHARSET UTF8MB4), 4, 0)) AS `date`,
         `piano`.`monthly_node_billing_info`.`node_uuid` AS `node_uuid`,
         `piano`.`monthly_node_billing_info`.`charge_cpu` AS `charge_cpu`,
         `piano`.`monthly_node_billing_info`.`charge_memory` AS `charge_memory`,
         `piano`.`monthly_node_billing_info`.`charge_nic` AS `charge_nic`,
         `piano`.`monthly_node_billing_info`.`uptime_ms` AS `uptime_ms`
     FROM
         `piano`.`monthly_node_billing_info`) `monthly`
GROUP BY `monthly`.`group_id` , `monthly`.`date` , `monthly`.`node_uuid`;

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
    SUM(`monthly`.`charge_traffic`) AS `charge_traffic`,
    SUM(`monthly`.`traffic_kb`) AS `traffic_kb`
FROM
    (SELECT
         `piano`.`monthly_server_billing_info`.`group_id` AS `group_id`,
         CONCAT(LPAD(CAST(FLOOR((CAST(REPLACE(`piano`.`monthly_server_billing_info`.`date`, '-', '')
                                     AS UNSIGNED) / 100))
                         AS CHAR (4) CHARSET UTF8MB4), 4, 0)) AS `date`,
         `piano`.`monthly_server_billing_info`.`server_uuid` AS `server_uuid`,
         `piano`.`monthly_server_billing_info`.`charge_traffic` AS `charge_traffic`,
         `piano`.`monthly_server_billing_info`.`traffic_kb` AS `traffic_kb`
     FROM
         `piano`.`monthly_server_billing_info`) `monthly`
GROUP BY `monthly`.`group_id` , `monthly`.`date` , `monthly`.`server_uuid`;

/* yearly_volume_billing_info */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`yearly_volume_billing_info` AS
SELECT
    `monthly`.`group_id` AS `group_id`,
    `monthly`.`date` AS `date`,
    `monthly`.`volume_uuid` AS `volume_uuid`,
    SUM(`monthly`.`charge_ssd`) AS `charge_ssd`,
    SUM(`monthly`.`charge_hdd`) AS `charge_hdd`
FROM
    (SELECT
         `piano`.`monthly_volume_billing_info`.`group_id` AS `group_id`,
         CONCAT(LPAD(CAST(FLOOR((CAST(REPLACE(`piano`.`monthly_volume_billing_info`.`date`, '-', '')
                                     AS UNSIGNED) / 100))
                         AS CHAR (4) CHARSET UTF8MB4), 4, 0)) AS `date`,
         `piano`.`monthly_volume_billing_info`.`volume_uuid` AS `volume_uuid`,
         `piano`.`monthly_volume_billing_info`.`charge_ssd` AS `charge_ssd`,
         `piano`.`monthly_volume_billing_info`.`charge_hdd` AS `charge_hdd`
     FROM
         `piano`.`monthly_volume_billing_info`) `monthly`
GROUP BY `monthly`.`group_id` , `monthly`.`date` , `monthly`.`volume_uuid`;
