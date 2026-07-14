resource "aws_db_subnet_group" "this" {
  name       = "openpayment-${var.environment}"
  subnet_ids = var.subnet_ids

  tags = {
    Name        = "openpayment-${var.environment}-db-subnet-group"
    Environment = var.environment
  }
}

resource "aws_security_group" "rds" {
  name        = "openpayment-${var.environment}-rds-sg"
  description = "Security group for RDS Postgres"
  vpc_id      = var.vpc_id

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    cidr_blocks     = [local.vpc_cidr]
    description     = "Postgres from VPC"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "openpayment-${var.environment}-rds-sg"
    Environment = var.environment
  }
}

resource "aws_db_instance" "this" {
  identifier     = "openpayment-${var.environment}"
  engine         = "postgres"
  engine_version = "16"
  instance_class = "db.r6g.large"

  db_name  = "paymentdb"
  username = "openpayment"
  password = var.db_password

  multi_az               = var.environment == "production" ? true : false
  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = [aws_security_group.rds.id]

  backup_retention_period = 30
  backup_window           = "03:00-04:00"
  maintenance_window      = "sun:05:00-sun:06:00"

  deletion_protection = true
  skip_final_snapshot = var.environment != "production"
  storage_encrypted   = true

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
