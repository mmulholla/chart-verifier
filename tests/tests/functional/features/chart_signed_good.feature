Feature: Chart verification
    Partners or redhat or community can verify their signed charts by running the
    chart verifier against an error free chart.

    Examples:
        | image_type |
        | docker |
        | tarball |
        | podman |

    Scenario Outline: A chart provider verifies their signed chart using the chart verifier
        Given I would like to use the <type> profile
        Given I will provide a <location> of a <helm_chart>
        Given I will provide a <location> of an expected <report_info>
        Given I will use the chart verifier <image_type> image
        Given I will provide a <location> of a <public_key> to verify the signature
        When I run the chart-verifier verify command against the signed chart to generate a report
        Then I should see the report-info from the report for the signed chart matching the expected report-info

        Examples:
            | type      | location                           | helm_chart               | report_info                | public_key                  |
            | partner   | tests/charts/psql-service/0.1.11/  | psql-service-0.1.11.tgz  | partner-report-info.json   | psql-service-0.1.11.tgz.key |
            | redhat    | tests/charts/psql-service/0.1.11/  | psql-service-0.1.11.tgz  | redhat-report-info.json    | psql-service-0.1.11.tgz.key |


