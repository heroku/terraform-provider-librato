package librato

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/heroku/go-librato/librato"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"librato": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderConfigureUsesDefaultBaseURLWhenNotSpecified(t *testing.T) {
	p := Provider().(*schema.Provider)
	d := schema.TestResourceDataRaw(t, p.Schema, nil)

	client, err := providerConfigure(d)
	if err != nil {
		t.Fatal(err)
	}

	c := client.(*librato.Client)

	if c.BaseURL.String() != "https://metrics-api.librato.com/v1/" {
		t.Fatalf("want url: https://metrics-api.librato.com/v1/, got %q", c.BaseURL.String())
	}
}

func TestProviderConfigureUsesCorrectURLEmailToken(t *testing.T) {
	p := Provider().(*schema.Provider)
	d := schema.TestResourceDataRaw(t, p.Schema, nil)
	d.Set("url", "https://some-url.com/v1/")
	d.Set("email", "foo@example.com")
	d.Set("token", "the-token")

	client, err := providerConfigure(d)
	if err != nil {
		t.Fatal(err)
	}

	c := client.(*librato.Client)

	if c.BaseURL.String() != "https://some-url.com/v1/" {
		t.Fatalf("want url: https://some-url.com/v1/, got %q", c.BaseURL.String())
	}

	if c.Email != "foo@example.com" {
		t.Fatalf("want email: foo@example.com, got %q", c.Email)
	}

	if c.Token != "the-token" {
		t.Fatalf("want token: the-token, got %q", c.Email)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("LIBRATO_EMAIL"); v == "" {
		t.Fatal("LIBRATO_EMAIL must be set for acceptance tests")
	}

	if v := os.Getenv("LIBRATO_TOKEN"); v == "" {
		t.Fatal("LIBRATO_TOKEN must be set for acceptance tests")
	}
}
