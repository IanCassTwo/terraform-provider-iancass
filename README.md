# terraform-provider-iancass
This provider implements the following features:-
* Application Load Balancer activations

## Authentication
    provider "iancass" {
        edgerc = "~/.edgerc"
        section = "default"
    }

## resource iancass_alb_activation
    resource "iancass_alb_activation" "alb" {
        origin_id = "icass_test"
        network = "staging"
        version = 3
    }

Note, this provider is to be used at your own risk. It's not an official Akamai release and shall not be supported in any way.
