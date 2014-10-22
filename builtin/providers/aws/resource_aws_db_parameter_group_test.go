package aws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/goamz/rds"
)

func TestAccAWSDBParameterGroup(t *testing.T) {
	var v rds.DBParameterGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDBParameterGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAWSDBParameterGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBParameterGroupExists("aws_db_parameter_group.bar", &v),
					testAccCheckAWSDBParameterGroupAttributes(&v),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "name", "parameter-group-test-terraform"),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "family", "mysql5.6"),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "description", "Test parameter group for terraform"),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "parameter.0.name", "character_set_client"),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "parameter.0.value", "utf8"),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "parameter.0.apply_method", "immediate"),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "parameter.1.name", "character_set_results"),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "parameter.1.value", "utf8"),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "parameter.1.apply_method", "immediate"),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "parameter.2.name", "character_set_server"),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "parameter.2.value", "utf8"),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "parameter.2.apply_method", "immediate"),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "parameter.3.name", "collation_connection"),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "parameter.3.value", "utf8_unicode_ci"),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "parameter.3.apply_method", "immediate"),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "parameter.4.name", "collation_server"),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "parameter.4.value", "utf8_unicode_ci"),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "parameter.4.apply_method", "immediate"),
				),
			},
		},
	})
}

func TestAccAWSDBParameterGroupOnly(t *testing.T) {
			var v rds.DBParameterGroup

			resource.Test(t, resource.TestCase{
				PreCheck:     func() { testAccPreCheck(t) },
				Providers:    testAccProviders,
				CheckDestroy: testAccCheckAWSDBParameterGroupDestroy,
				Steps: []resource.TestStep{
				resource.TestStep{
				Config: testAccAWSDBParameterGroupOnlyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDBParameterGroupExists("aws_db_parameter_group.bar", &v),
					testAccCheckAWSDBParameterGroupAttributes(&v),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "name", "parameter-group-test-terraform"),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "family", "mysql5.6"),
					resource.TestCheckResourceAttr(
						"aws_db_parameter_group.bar", "description", "Test parameter group for terraform"),
				),
			},
		},
	})
}

func testAccCheckAWSDBParameterGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.rdsconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_db_parameter_group" {
			continue
		}

		// Try to find the Group
		resp, err := conn.DescribeDBParameterGroups(
			&rds.DescribeDBParameterGroups{
				DBParameterGroupName: rs.Primary.ID,
			})

		if err == nil {
			if len(resp.DBParameterGroups) != 0 &&
				resp.DBParameterGroups[0].DBParameterGroupName == rs.Primary.ID {
				return fmt.Errorf("DB Parameter Group still exists")
			}
		}

		// Verify the error
		newerr, ok := err.(*rds.Error)
		if !ok {
			return err
		}
		if newerr.Code != "InvalidDBParameterGroup.NotFound" {
			return err
		}
	}

	return nil
}

func testAccCheckAWSDBParameterGroupAttributes(v *rds.DBParameterGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if v.DBParameterGroupName != "parameter-group-test-terraform" {
			return fmt.Errorf("bad name: %#v", v.DBParameterGroupName)
		}

		if v.DBParameterGroupFamily != "mysql5.6" {
			return fmt.Errorf("bad family: %#v", v.DBParameterGroupFamily)
		}

		if v.Description != "Test parameter group for terraform" {
			return fmt.Errorf("bad description: %#v", v.Description)
		}

		return nil
	}
}

func testAccCheckAWSDBParameterGroupExists(n string, v *rds.DBParameterGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB Parameter Group ID is set")
		}

		conn := testAccProvider.rdsconn

		opts := rds.DescribeDBParameterGroups{
			DBParameterGroupName: rs.Primary.ID,
		}

		resp, err := conn.DescribeDBParameterGroups(&opts)

		if err != nil {
			return err
		}

		if len(resp.DBParameterGroups) != 1 ||
			resp.DBParameterGroups[0].DBParameterGroupName != rs.Primary.ID {
			return fmt.Errorf("DB Parameter Group not found")
		}

		*v = resp.DBParameterGroups[0]

		return nil
	}
}

const testAccAWSDBParameterGroupConfig = `
resource "aws_db_parameter_group" "bar" {
	name = "parameter-group-test-terraform"
	family = "mysql5.6"
	description = "Test parameter group for terraform"
	parameter {
	  name = "character_set_server"
	  value = "utf8"
	  apply_method = "immediate"
	}
	parameter {
	  name = "character_set_client"
	  value = "utf8"
	  apply_method = "immediate"
	}
	parameter{
	  name = "character_set_results"
	  value = "utf8"
	  apply_method = "immediate"
	}
	parameter {
	  name = "collation_server"
	  value = "utf8_unicode_ci"
	  apply_method = "immediate"
	}
	parameter {
	  name = "collation_connection"
	  value = "utf8_unicode_ci"
	  apply_method = "immediate"
	}
}
`

const testAccAWSDBParameterGroupOnlyConfig = `
resource "aws_db_parameter_group" "bar" {
	name = "parameter-group-test-terraform"
	family = "mysql5.6"
	description = "Test parameter group for terraform"
}
`