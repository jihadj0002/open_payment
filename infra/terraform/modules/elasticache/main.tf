resource "aws_elasticache_subnet_group" "this" {
  name       = "openpayment-${var.environment}"
  subnet_ids = var.subnet_ids

  tags = {
    Name        = "openpayment-${var.environment}-redis-subnet-group"
    Environment = var.environment
  }
}

resource "aws_security_group" "redis" {
  name        = "openpayment-${var.environment}-redis-sg"
  description = "Security group for ElastiCache Redis"
  vpc_id      = var.vpc_id

  ingress {
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    cidr_blocks     = [local.vpc_cidr]
    description     = "Redis from VPC"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "openpayment-${var.environment}-redis-sg"
    Environment = var.environment
  }
}

resource "aws_elasticache_replication_group" "this" {
  replication_group_id          = "openpayment-${var.environment}"
  description                   = "Redis ${var.environment}"
  engine                        = "redis"
  engine_version                = "7.0"
  node_type                     = "cache.r6g.large"
  num_cache_clusters            = 1
  port                          = 6379
  parameter_group_name          = "default.redis7"
  subnet_group_name             = aws_elasticache_subnet_group.this.name
  security_group_ids            = [aws_security_group.redis.id]
  automatic_failover_enabled    = false
  multi_az_enabled              = false
  at_rest_encryption_enabled    = true
  transit_encryption_enabled    = true
  snapshot_retention_limit      = 7
  snapshot_window               = "03:00-04:00"
  maintenance_window            = "sun:06:00-sun:07:00"

  tags = {
    Name        = "openpayment-${var.environment}"
    Environment = var.environment
  }
}

data "aws_vpc" "selected" {
  id = var.vpc_id
}

locals {
  vpc_cidr = data.aws_vpc.selected.cidr_block
}
