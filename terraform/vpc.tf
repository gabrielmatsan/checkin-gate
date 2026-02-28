resource "aws_vpc" "main-vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name        = "main-vpc"
    Environment = var.environment
    IaC         = "true"
  }
}

resource "aws_subnet" "public-a" {
  vpc_id = aws_vpc.main-vpc.id

  cidr_block        = "10.0.1.0/24"
  availability_zone = "us-east-1a"


  tags = {
    Name        = "public-a"
    Environment = var.environment
    IaC         = "true"
  }
}

resource "aws_subnet" "public-b" {
  vpc_id = aws_vpc.main-vpc.id

  cidr_block        = "10.0.2.0/24"
  availability_zone = "us-east-1b"

  tags = {
    Name        = "public-b"
    Environment = var.environment
    IaC         = "true"
  }
}


resource "aws_subnet" "private-a" {
  vpc_id = aws_vpc.main-vpc.id

  cidr_block        = "10.0.3.0/24"
  availability_zone = "us-east-1a"

  tags = {
    Name        = "private-a"
    Environment = var.environment
    IaC         = "true"
  }
}

resource "aws_subnet" "private-b" {
  vpc_id = aws_vpc.main-vpc.id

  cidr_block        = "10.0.4.0/24"
  availability_zone = "us-east-1b"

  tags = {
    Name        = "private-b"
    Environment = var.environment
    IaC         = "true"
  }
}


resource "aws_internet_gateway" "main-igw" {
  vpc_id = aws_vpc.main-vpc.id

  tags = {
    Name        = "main-igw"
    Environment = var.environment
    IaC         = "true"
  }
}

# NAT
resource "aws_eip" "eip" {
  domain = "vpc"

  tags = {
    Name        = "main-eip"
    Environment = var.environment
    IaC         = "true"
  }
}


resource "aws_nat_gateway" "nat" {
  allocation_id = aws_eip.eip.id
  subnet_id     = aws_subnet.public-a.id

  tags = {
    Name        = "main-nat"
    Environment = var.environment
    IaC         = "true"
  }
}


resource "aws_route_table" "public-route-table" {
  vpc_id = aws_vpc.main-vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main-igw.id
  }


  tags = {
    Name        = "public-route-table"
    Environment = var.environment
    IaC         = "true"
  }

}


resource "aws_route_table_association" "public-a-association" {
  subnet_id      = aws_subnet.public-a.id
  route_table_id = aws_route_table.public-route-table.id
}


resource "aws_route_table_association" "public-b-association" {
  subnet_id      = aws_subnet.public-b.id
  route_table_id = aws_route_table.public-route-table.id
}

resource "aws_route_table" "private-route-table" {
  vpc_id = aws_vpc.main-vpc.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.nat.id
  }

  tags = {
    Name        = "private-route-table"
    Environment = var.environment
    IaC         = "true"
  }
}


resource "aws_route_table_association" "private-a-association" {
  subnet_id      = aws_subnet.private-a.id
  route_table_id = aws_route_table.private-route-table.id
}

resource "aws_route_table_association" "private-b-association" {
  subnet_id      = aws_subnet.private-b.id
  route_table_id = aws_route_table.private-route-table.id
}

output "vpc_id" {
  value = aws_vpc.main-vpc.id
}
output "public_subnet_ids" {
  value = [aws_subnet.public-a.id, aws_subnet.public-b.id]
}
output "private_subnet_ids" {
  value = [aws_subnet.private-a.id, aws_subnet.private-b.id]
}