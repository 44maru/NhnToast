package constants

const FLOW_TYPE_CREATE_INSTANCE = "create-instance"
const FLOW_TYPE_DELETE_INSTANCE = "delete-instance"
const FLOW_TYPE_LIST_INSTANCE = "list-instance"
const FLOW_TYPE_CREATE_FLOATINGIP = "create-floatingip"
const FLOW_TYPE_DELETE_FLOATINGIP = "delete-floatingip"
const FLOW_TYPE_LIST_FLOATINGIP = "list-floatingip"

const COMPUTE_ENDPOINT = "https://jp1-api-instance.infrastructure.cloud.toast.com"
const NETWORK_ENDPOINT = "https://jp1-api-network.infrastructure.cloud.toast.com"
const IMAGE_ENDPOINT = "https://jp1-api-image.infrastructure.cloud.toast.com"

const FLOATING_IP_URL = NETWORK_ENDPOINT + "/v2.0/floatingips"
const NETWORK_URL = NETWORK_ENDPOINT + "/v2.0/networks"
const SUBNET_URL = NETWORK_ENDPOINT + "/v2.0/subnets"
const PORT_URL = NETWORK_ENDPOINT + "/v2.0/ports"
const IMAGE_URL = IMAGE_ENDPOINT + "/v2/images"

const PUBLIC_NETWORK_ID = "117fa565-c8eb-4e58-a420-c5146e516341"
const DEFAULT_SUBNET_NAME = "Default Network"
const PUBLIC_NETWORK_NAME = "Public network for Toast JP"
