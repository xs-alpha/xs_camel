create table `accessory_price_query`(
`accessory_price_query_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
`vehicle_type` varchar(100) NOT NULL DEFAULT '' COMMENT '车型',
`vehicle_series` varchar(100) NOT NULL DEFAULT '' COMMENT '车系',
`accessory_type` varchar(100) NOT NULL DEFAULT '' COMMENT '类型',
`material_no` varchar(40) NOT NULL DEFAULT '' COMMENT '物料代码',
`accessory_title` varchar(100) NOT NULL DEFAULT '' COMMENT '配件名称',
`price`  decimal(13,2) null COMMENT '价格',
`create_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`update_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (`accessory_price_query_id`)
)ENGINE=InnoDB   COMMENT='配件价格查询表';



11  22  33  444 
22  33  44  55

11
22
33
44
5555
66