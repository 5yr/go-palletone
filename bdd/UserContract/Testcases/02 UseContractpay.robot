*** Settings ***
Resource          ../../commonlib/pubVariables.robot
Resource          ../../commonlib/pubFuncs.robot
Resource          ../../commonlib/setups.robot
Library           BuiltIn

*** Test Cases ***
InstallContractpayTpl
    Given Unlock token holder succeed
    ${reqId} =    When User installs contract template
    And wait for transaction being packaged
    Then Wait for unit abount contract to be confirmed by unit height    ${reqId}

DeployContract
    Given Unlock token holder succeed
    ${reqId} =    When User deploys contract
    And wait for transaction being packaged
    Then Wait for unit abount contract to be confirmed by unit height    ${reqId}

PutStatus
    Given Unlock token holder succeed
    ${reqId} =    When User put status into contractpay
    And wait for transaction being packaged
    And Wait for unit abount contract to be confirmed by unit height    ${reqId}
    Then Get status from contractpay

Payout
    Given Unlock token holder succeed
    ${reqId} =    When User transfer PTN to contractpay
    And wait for transaction being packaged
    And Wait for unit abount contract to be confirmed by unit height    ${reqId}
    And Query contract balance
    ${newAddr}    ${reqId}=    And Use contractpay to transfer PTN to user2
    And wait for transaction being packaged
    ${reqId} =    And Wait for unit abount contract to be confirmed by unit height    ${reqId}
    Then Query user2 balance

*** Keywords ***
Unlock token holder succeed
    unlockAccount    ${tokenHolder}

User installs contract template
    ${respJson}=    installContractTpl    ${tokenHolder}    ${tokenHolder}    100    100    jury06
    ...    github.com/palletone/go-palletone/contracts/example/go/contractpay    example
    ${result}=    Get From Dictionary    ${respJson}    result
    ${reqId}=    Get From Dictionary    ${result}    reqId
    ${tplId}=    Get From Dictionary    ${result}    tplId
    Set Global Variable    ${gTplId}    ${tplId}
    [Return]    ${reqId}

User deploys contract
    ${args}=    Create List    deploy
    ${respJson}=    deployContract    ${tokenHolder}    ${tokenHolder}    1000    10    ${gTplId}
    ...    ${args}
    ${result}=    Get From Dictionary    ${respJson}    result
    ${reqId}=    Get From Dictionary    ${result}    reqId
    ${contractId}=    Get From Dictionary    ${result}    ContractId
    Set Global Variable    ${gContractId}    ${contractId}
    [Return]    ${reqId}

User put status into contractpay
    ${args}=    Create List    put    Hello
    ${respJson}=    invokeContract    ${tokenHolder}    ${tokenHolder}    100    1    ${gContractId}
    ...    ${args}
    ${result}=    Get From Dictionary    ${respJson}    result
    ${reqId}=    Get From Dictionary    ${result}    reqId
    [Return]    ${reqId}

Get status from contractpay
    ${args}=    Create List    get    Hello
    ${respJson}=    queryContract    ${gContractId}    ${args}
    Log    ${respJson}

User transfer PTN to contractpay
    transferPtnTo    ${gContractId}    10000
    wait for transaction being packaged

Query contract balance
    ${amount}=    getBalance    ${gContractId}
    Should Be Equal    ${amount}    10000
    Log    ${amount}

Use contractpay to transfer PTN to user2
    # create account
    ${newAddr}=    newAccount
    ${args}=    Create List    payout    ${newAddr}    PTN    100
    ${respJson}=    invokeContract    ${tokenHolder}    ${tokenHolder}    100    1    ${gContractId}
    ...    ${args}
    ${reqId}=    Get From Dictionary    ${respJson}    reqId
    [Return]    ${newAddr}    ${reqId}

Query user2 balance
    [Arguments]    ${addr}
    ${args}=    Create List    balance    ${addr}
    ${respJson}=    invokeContract    ${tokenHolder}    ${tokenHolder}    100    1    ${gContractId}
    ...    ${args}
    Log    ${respJson}
    # [Return]    ${reqId}
