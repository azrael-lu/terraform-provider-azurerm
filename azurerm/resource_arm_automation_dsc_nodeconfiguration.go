package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAutomationDscNodeConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationDscNodeConfigurationCreateUpdate,
		Read:   resourceArmAutomationDscNodeConfigurationRead,
		Update: resourceArmAutomationDscNodeConfigurationCreateUpdate,
		Delete: resourceArmAutomationDscNodeConfigurationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceArmAutomationDscNodeConfigurationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationDscNodeConfigurationClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Automation Dsc Node Configuration creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	accName := d.Get("account_name").(string)
	content := d.Get("content").(string)

	s := strings.Split(name, ".")

	parameters := automation.DscNodeConfigurationCreateOrUpdateParameters{
		Source: &automation.ContentSource{
			Type:  automation.EmbeddedContent,
			Value: &content,
		},
		Configuration: &automation.DscConfigurationAssociationProperty{
			Name: &s[0],
		},
		Name: &name,
	}

	_, err := client.CreateOrUpdate(ctx, resGroup, accName, name, parameters)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, accName, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Automation Dsc Node Configuration '%s' (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAutomationDscNodeConfigurationRead(d, meta)
}

func resourceArmAutomationDscNodeConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationDscNodeConfigurationClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["nodeConfigurations"]

	resp, err := client.Get(ctx, resGroup, accName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM Automation Dsc Node Configuration '%s': %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("account_name", accName)

	return nil
}

func resourceArmAutomationDscNodeConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationDscNodeConfigurationClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["nodeConfigurations"]

	resp, err := client.Delete(ctx, resGroup, accName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error issuing AzureRM delete request for Automation Dsc Node Configuration '%s': %+v", name, err)
	}

	return nil
}
