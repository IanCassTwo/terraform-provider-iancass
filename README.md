# terraform-provider-iancass

WARNING: Please use this as example code. This is not production ready and may not work as expected. Use at your own risk.

This provider implements the following features:-
* Application Load Balancer activations
* CPS DV certificate
* CPS Third Part certificate (untested)
* Firewall Rules Notification (untested)
* Siteshield (datasource only) (untested)

## Authentication
    provider "iancass" {
        edgerc = "~/.edgerc"
        config_section = "default"
    }

## resource iancass_alb_activation
    resource "iancass_alb_activation" "alb" {
        origin_id = "icass_test"
        network = "staging"
        version = 3
    }

## resource iancass_cps_dv_enrollment
        resource "iancass_cps_dv_enrollment" "cert" {
            organization {
                    name = "Example Name"
                    phone = "01212323323"
                    addresslineone = "123 Example Street"
                    addresslinetwo = "Exampleton"
                    city = "Exampledon"
                    region = "Exampleshire"
                    postalcode = "WI553BL"
                    countrycode = "GB"
            }
            admincontact {
                    firstname = "Joe"
                    lastname = "Public"
                    title = "Mr"
                    organization = "Akamai"
                    email = "joepublic@example.com"
                    phone = "012312321313"
                    addresslineone = "123 Example Street"
                    addresslinetwo = "Exampleton"
                    city = "Exampledon"
                    region = "Exampleshire"
                    postalcode = "WI553BL"
                    countrycode = "GB"
            }
            techcontact {
                    firstname = "Joe"
                    lastname = "Public"
                    title = "Mr"
                    organization = "Akamai"
                    email = "joepublic@example.com"
                    phone = "012312321313"
                    addresslineone = "123 Example Street"
                    addresslinetwo = "Exampleton"
                    city = "Exampledon"
                    region = "Exampleshire"
                    postalcode = "WI553BL"
                    countrycode = "GB"
            }
            snionly = true
            securenetwork = "standard-tls"
            commonname = "test1.example.com"
            alternativenames = [ "test2.example.com" ]
            contract = "1-XXXXXX"
    }

    output "dns-challenges" {
            value = iancass_cps_dv_enrollment.cert.dnschallenges
    }

## iancass_cps_dv_validation
    resource "iancass_cps_dv_validation" "validation" {
            certificateid = iancass_cps_dv_enrollment.cert.id
            depends_on = [
                    akamai_dns_record.validation,
            ]
    }

WARNING: This provider is to be used at your own risk. It's not an official Akamai release and shall not be supported in any way.
