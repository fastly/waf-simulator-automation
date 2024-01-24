# waf-simulator-automation

This project contains tooling, written in Go, that is designed to facilitate automated testing against Fastly's WAF simulator. In particular, we’ve incorporated a CI/CD pipeline that uses Github action workflows to run tests on every code change in the main branch.

## Dependencies 

- [Fastly NGWAF](https://www.fastly.com/products/web-application-api-protection)
- [Go](https://go.dev/doc/install)
- [Terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)

## Test Configurations

Tests are written in yaml format and located in the test/rules directory. The yaml files serve as a structured way to define and organize test cases. Each test file contains a list of tests and each test contains the following fields.

- **name**: (required) A unique identifier name for the test case.  
- **site**: (required) Identifying name of the site in Fastly NGWAF to test against.  
- **rule_id**: (optional) ID of rule you are testing against.  
- **description**: (optional) Details about what the test is designed to check.  
- **type**: (optional) True positive, false negative, false positive, true negative.  
- **request**: (required) HTTP request that will be sent as part of the test.  
- **response**: (required) The expected response for the test.  
- **expect**: (required) This section outlines the expected outcome of the test.  
  - **waf_response:** The expected response code from the WAF.  
  - **signals:** A list of the signaled data to be returned by the test. Each signal contains several values and should be omitted if empty.  
    - **type:** Signal ID (a.k.a signal type)
    - **location:** Location of signaled value (i.e. QUERYSTRING, USERAGENT)
    - **name:** The name assigned to the signal
    - **value:** The specific value that triggered the signal
    - **detector:** The identifier of the detector that generated the signal
    - **redaction:** A binary indicator (1 or 0) signifying whether the signal's value has been redacted.

## Getting Started

Follow the steps below:

1. Clone the repository [https://github.com/fastly/waf-simulator-automation](https://github.com/fastly/waf-simulator-automation)
2. Create an NGWAF API key
    - Sign into the NGWAF console at [https://dashboard.signalsciences.net/login](https://dashboard.signalsciences.net/login)
    - On the **My Profile** tab, under **API Access Tokens** , select **Add API access token**.
    - Type in a name and select **Create API access token.**
3. Set your Fastly NGWAF credentials as environment variables.

    ```bash
    export SIGSCI_EMAIL='your-email'
    export SIGSCI_TOKEN='your-token'
    export SIGSCI_CORP='your-corp-id'
    ```
4. Install Terraform if not already installed with the steps described [here](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli).
5. From the project directory change into the **terraform** directory and run the following commands.

    ```bash
    terraform init
    terraform plan -out ngwaf.plan
    terraform apply ngwaf.plan
    ```

6. After running apply, take note of the output values of **sensitive\_account\_api\_rule\_id** and **invalid\_host\_header\_rule\_id**.
7. Open **tests/rules/app1.example.com/sensitive-account-api.yaml** and replace all occurrences of **65a190f3e3148001dc71a5ca** with the **sensitive\_account\_api\_rule\_id** value from the terraform output.
8. Open **tests/rules/app2.example.com/invalid-host-headers.yaml** file and replace all occurrences of **65a190f40f6eb201dc0fdd81** with the **invalid\_host\_header\_rule\_id** value from the terraform output.
9. Once the test files have been updated you can run the WAF simulator tests to verify the WAF rules are working correctly.
10. Install Go if not already installed using the steps described [here](https://go.dev/doc/install).
11. Change back to the project's root directory and run the following command.

    ```bash
    go run tests/main.go
    ```

12. If you didn't receive any failures, the tests passed. If you see failures, use the logs to troubleshoot and resolve the issues.
13. Create a new repository on GitHub with steps described [here](https://docs.github.com/en/repositories/creating-and-managing-repositories/quickstart-for-repositories).
14. Change the remote URL to your new repository.
    - In your terminal or command prompt, navigate to the cloned repository's directory.
    - Use the **git remote** command to change the remote URL to your new repository. This points your local repository to the new GitHub repository.

    ```bash 
    git remote set-url origin https://github.com/yourusername/new-repository.git
    ```

15. Add **SIGSCI\_EMAIL** , **SIGSCI\_CORP** , **SIGSCI\_TOKEN** to [GitHub secrets](https://docs.github.com/en/actions/security-guides/using-secrets-in-github-actions#creating-secrets-for-a-repository)
16. Now, push the code to your new repository using the **git push** command.
    ```bash
    git push origin main
    ```
17. In the workflow file [.github/workflows/tests.yaml](.github/workflows/tests.yaml#L5), change the branch name from main-branch to main.
    ```bash
    git add .github/workflows/tests.yaml
    git commit -m "update workflow"
    ``` 
18. After pushing, check your repository on GitHub to ensure the test workflow is running.
19. In your repository, locate the **Actions** tab near the top of the page. This tab shows you a list of workflow runs associated with your repository. You'll see a list of recent workflow runs. Each run is associated with a commit or event that triggered it (like a push to the main branch).
20. If the workflow succeeded, your WAF rules are working as expected.
21. If there are failures, use the logs to troubleshoot and resolve issues. After making corrections, commit and push your changes again to trigger the workflow.
