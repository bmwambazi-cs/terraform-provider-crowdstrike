terraform {
  required_providers {
    crowdstrike = {
      source = "registry.terraform.io/crowdstrike/crowdstrike"
    }
  }
}

provider "crowdstrike" {
  cloud = "us-2"
}


resource "crowdstrike_cloud_compliance_custom_framework" "example" {
  name        = "example-framework"
  description = "An example framework created with Terraform"
  sections = {
    "section-1" = { // immutable unique key
      name = "Section 1"
      controls = {
        "control-1a" = { // immutable unique key
          name        = "Control 1a"
          description = "This is the first control"
          rules       = ["id1", "id2", "id3"]
        }
        "control-1b" = {
          name        = "Control 1b"
          description = "This is another control in section 1"
          rules       = ["id4", "id5"]
        }
      }
    }
    "section-2" = {
      name = "Section 2"
      controls = {
        "control-2" = {
          name        = "Control 2"
          description = "This is the second control"
          rules       = []
        }
      }
    }
  }
}

// Clone an existing framework by providing its ID.
// Sections and controls are inherited from the parent framework.
resource "crowdstrike_cloud_compliance_custom_framework" "cloned" {
  name                = "cloned-framework"
  description         = "A framework cloned from an existing one"
  parent_framework_id = "7c86a274-c04b-4292-9f03-dafae42bde97"
}

// Clone a framework and override specific sections/controls.
// Cloned sections not mentioned here are preserved. Config sections are merged
// on top: matching section keys merge controls within, new keys are added.
resource "crowdstrike_cloud_compliance_custom_framework" "cloned_with_overrides" {
  name                = "cloned-with-overrides"
  description         = "A cloned framework with custom section overrides"
  parent_framework_id = "7c86a274-c04b-4292-9f03-dafae42bde97"
  sections = {
    "custom-section" = {
      name = "Custom Section"
      controls = {
        "custom-control" = {
          name        = "Custom Control"
          description = "A control added on top of the cloned framework"
          rules       = []
        }
      }
    }
  }
}

output "cloud_compliance_custom_framework" {
  value = crowdstrike_cloud_compliance_custom_framework.example
}

output "cloud_compliance_custom_framework_cloned" {
  value = crowdstrike_cloud_compliance_custom_framework.cloned
}

output "cloud_compliance_custom_framework_cloned_with_overrides" {
  value = crowdstrike_cloud_compliance_custom_framework.cloned_with_overrides
}
