# InsightFinder Project with Holiday Settings

This example demonstrates how to configure holiday settings for an InsightFinder project using the Terraform provider.

## Holiday Settings

Holiday settings allow you to define specific dates that should be treated as holidays in your project. This is useful for:
- Adjusting anomaly detection during known holiday periods
- Accounting for expected traffic/behavior changes during holidays
- Improving model accuracy by factoring in holiday patterns

### Format

Each holiday requires three fields:
- **name**: A unique identifier for the holiday (string)
- **start_date**: The start date in MM-DD format (e.g., "12-25")
- **end_date**: The end date in MM-DD format (e.g., "12-26")

### Example Usage

```hcl
resource "insightfinder_project" "example_with_holidays" {
  project_name = "example-project"
  system_name  = "example-system"

  project_creation_config {
    data_type          = "Log"
    instance_type      = "PrivateCloud"
    project_cloud_type = "PrivateCloud"
  }

  holiday_settings = [
    {
      name       = "christmas"
      start_date = "12-25"
      end_date   = "12-26"
    },
    {
      name       = "new_year"
      start_date = "01-01"
      end_date   = "01-01"
    }
  ]
}
```

## Usage

1. Copy this example to a new directory
2. Update the variables in `terraform.tfvars` or set environment variables:
   ```bash
   export TF_VAR_license_key="your-license-key"
   export TF_VAR_user_name="your-username"
   export TF_VAR_server_url="https://stg.insightfinder.com"
   ```

3. Initialize Terraform:
   ```bash
   terraform init
   ```

4. Review the plan:
   ```bash
   terraform plan
   ```

5. Apply the configuration:
   ```bash
   terraform apply
   ```

## Notes

- Date format must be MM-DD (e.g., "12-25" for December 25th)
- Holiday names must be unique within a project
- Holidays are stored and retrieved via the InsightFinder API
- When updating holidays, the provider will:
  - Delete holidays that are removed from the configuration
  - Update holidays that have changed dates
  - Add new holidays that don't exist yet

## API Endpoints Used

- **GET** `/api/external/v1/holiday` - Retrieve holidays
- **POST** `/api/external/v1/holiday` - Create a holiday
- **DELETE** `/api/external/v1/holiday` - Delete holidays
