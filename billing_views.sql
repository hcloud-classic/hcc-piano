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
    `piano`.`daily_info`.`charge_node` AS `charge_node`,
    `piano`.`daily_info`.`charge_server` AS `charge_server`,
    `piano`.`daily_info`.`charge_network` AS `charge_network`,
    `piano`.`daily_info`.`charge_volume` AS `charge_volume`
FROM
    `piano`.`daily_info`

/* monthly_bill */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`monthly_bill` AS
SELECT
    `daily`.`date` AS `date`,
    `daily`.`group_id` AS `group_id`,
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
         `piano`.`daily_bill`.`charge_node` AS `charge_node`,
         `piano`.`daily_bill`.`charge_server` AS `charge_server`,
         `piano`.`daily_bill`.`charge_network` AS `charge_network`,
         `piano`.`daily_bill`.`charge_volume` AS `charge_volume`
     FROM
         `piano`.`daily_bill`) `daily`
GROUP BY `daily`.`date` , `daily`.`group_id`

/* yearly_bill */
CREATE
    ALGORITHM = UNDEFINED
    DEFINER = `root`@`%`
    SQL SECURITY DEFINER
    VIEW `piano`.`yearly_bill` AS
SELECT
    `monthly`.`date` AS `date`,
    `monthly`.`group_id` AS `group_id`,
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
         `piano`.`monthly_bill`.`charge_node` AS `charge_node`,
         `piano`.`monthly_bill`.`charge_server` AS `charge_server`,
         `piano`.`monthly_bill`.`charge_network` AS `charge_network`,
         `piano`.`monthly_bill`.`charge_volume` AS `charge_volume`
     FROM
         `piano`.`monthly_bill`) `monthly`
GROUP BY `monthly`.`date` , `monthly`.`group_id`