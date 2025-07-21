variable "vpc_name" {
  description = "Name of the VPC"
  type        = string
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
}

variable "azs" {
  description = "List of availability zones"
  type        = list(string)
}

variable "private_subnets" {
  description = "List of private subnet CIDR blocks"
  type        = list(string)
}

variable "public_subnets" {
  description = "List of public subnet CIDR blocks"
  type        = list(string)
}

variable "enable_nat_gateway" {
  description = "Whether to enable NAT gateway"
  type        = bool
  default     = true
}

variable "single_nat_gateway" {
  description = "Whether to use a single NAT gateway"
  type        = bool
  default     = false
}

variable "enable_dns_hostnames" {
  description = "Whether to enable DNS hostnames"
  type        = bool
  default     = true
}

# Placeholder VPC (will be replaced with actual implementation)
output "vpc_id" {
  description = "ID of the VPC"
  value       = "vpc-placeholder"
}

output "private_subnets" {
  description = "List of private subnet IDs"
  value       = ["subnet-private-placeholder"]
}

output "public_subnets" {
  description = "List of public subnet IDs"
  value       = ["subnet-public-placeholder"]
}
