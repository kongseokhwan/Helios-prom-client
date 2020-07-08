package ovs_prom_ctx

const OVS_INTERFACE_RECEIVE_BYTES_TOTAL string = "ovs_interface_receive_bytes_total"
const OVS_INTERFACE_RECEIVE_CRC_TOTAL string = "ovs_interface_receive_crc_total"
const OVS_INTERFACE_RECEIVE_DROP_TOTAL string = "ovs_interface_receive_drop_total"
const OVS_INTERFACE_RECEIVE_ERRORS_TOTAL string = "ovs_interface_receive_errors_total"
const OVS_INTERFACE_RECEIVE_PACKETS_TOTAL string = "ovs_interface_receive_packets_total"
const OVS_INTERFACE_TRANSMIT_BYTES_TOTAL string = "ovs_interface_transmit_bytes_total"
const OVS_INTERFACE_TRANSMIT_COLLISIONS_TOTAL string = "ovs_interface_transmit_collisionss_total"
const OVS_INTERFACE_TRANSMIT_DROP_TOTAL string = "ovs_interface_transmit_drop_total"
const OVS_INTERFACE_TRANSMIT_ERRORS_TOTAL string = "ovs_interface_transmit_errors_total"
const OVS_INTERFACE_TRANSMIT_PACKETS_TOTAL string = "ovs_interface_transmit_packeets_total"

const OVS_FLOW_FLOW_BYTES_TOTAL string = "ovs_flow_flow_bytes_total"
const OVS_FLOW_FLOW_PACKETS_TOTAL string = "ovs_flow_flow_packets_total"

const ntopQueryWithRate string = "topk(%s, avg by (bridge, port)(rate(%s[%s])*8))" // rankSize, metric, duration
const countQuery string = "count(count by (bridge, port)(%s)"                      // metric
const avgbyQueryWithRate string = "avg by(bridge, port) (rate(%s[%s])*8)"          // metric, duration
