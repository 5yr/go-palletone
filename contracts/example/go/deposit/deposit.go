/*
	This file is part of go-palletone.
	go-palletone is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.
	go-palletone is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.
	You should have received a copy of the GNU General Public License
	along with go-palletone.  If not, see <http://www.gnu.org/licenses/>.
*/
/*
 * Copyright IBM Corp. All Rights Reserved.
 * @author PalletOne core developers <dev@pallet.one>
 * @date 2018
 */

//Package deposit implements some functions for deposit contract.
package deposit

import (
	"encoding/json"
	"fmt"
	"github.com/palletone/go-palletone/common/award"
	"github.com/palletone/go-palletone/contracts/shim"
	pb "github.com/palletone/go-palletone/core/vmContractPub/protos/peer"
	"github.com/palletone/go-palletone/dag/modules"
	"strconv"
	"strings"
	"time"
)

var (
	depositAmountsForJury      uint64
	depositAmountsForMediator  uint64
	depositAmountsForDeveloper uint64
	depositPeriod              int
	foundationAddress          string
)

type DepositChaincode struct {
}

func (d *DepositChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("*** DepositChaincode system contract init ***")
	depositPeriod, err := stub.GetSystemConfig("DepositPeriod")
	if err != nil {
		return shim.Error(err.Error())
	}
	day, _ := strconv.Atoi(depositPeriod)
	fmt.Println("加入保证金候选列表，需要持币在规定时间以上，规定时间为 = ", day)
	fmt.Println()
	foundationAddress, err = stub.GetSystemConfig("FoundationAddress")
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("foundationAddress = ", foundationAddress)
	fmt.Println()
	depositAmountsForMediatorStr, err := stub.GetSystemConfig("DepositAmountForMediator")
	if err != nil {
		return shim.Success([]byte("GetSystemConfig with DepositAmount error: "))
	}
	//转换
	depositAmountsForMediator, err = strconv.ParseUint(depositAmountsForMediatorStr, 10, 64)
	if err != nil {
		return shim.Success([]byte("String transform to uint64 error:"))
	}
	fmt.Println("需要的mediator保证金数量=", depositAmountsForMediator)
	fmt.Println()
	depositAmountsForJuryStr, err := stub.GetSystemConfig("DepositAmountForJury")
	if err != nil {
		return shim.Success([]byte("GetSystemConfig with DepositAmount error:"))
	}
	//转换
	depositAmountsForJury, err = strconv.ParseUint(depositAmountsForJuryStr, 10, 64)
	if err != nil {
		return shim.Success([]byte("String transform to uint64 error:"))
	}
	fmt.Println("需要的jury保证金数量=", depositAmountsForJury)
	fmt.Println()
	depositAmountsForDeveloperStr, err := stub.GetSystemConfig("DepositAmountForDeveloper")
	if err != nil {
		return shim.Success([]byte("GetSystemConfig with DepositAmount error:"))
	}
	//转换
	depositAmountsForDeveloper, err = strconv.ParseUint(depositAmountsForDeveloperStr, 10, 64)
	if err != nil {
		return shim.Success([]byte("String transform to uint64 error:"))
	}
	fmt.Println("需要的Developer保证金数量=", depositAmountsForDeveloper)
	fmt.Println()
	return shim.Success([]byte("ok"))
}

func (d *DepositChaincode) mediatorPayToDepositContract(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return mediatorPayToDepositContract(stub, args)
}

func (d *DepositChaincode) juryPayToDepositContract(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return juryPayToDepositContract(stub, args)
}

func (d *DepositChaincode) developerPayToDepositContract(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return developerPayToDepositContract(stub, args)
}
func (d *DepositChaincode) mediatorApplyCashback(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return mediatorApplyCashback(stub, args)
}
func (d *DepositChaincode) juryApplyCashback(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return juryApplyCashback(stub, args)
}
func (d *DepositChaincode) developerApplyCashback(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return developerApplyCashback(stub, args)
}

func (d *DepositChaincode) handleForMediatorApplyCashback(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return handleForMediatorApplyCashback(stub, args)
}

func (d *DepositChaincode) handleForJuryApplyCashback(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return handleForJuryApplyCashback(stub, args)
}

func (d *DepositChaincode) handleForDeveloperApplyCashback(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return handleForDeveloperApplyCashback(stub, args)
}

func (d *DepositChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()
	//fmt.Println(funcName)
	//for _, v := range args {
	//	fmt.Println(string(v))
	//}
	switch funcName {
	case "ApplyBecomeMediator":
		//申请成为Mediator
		return d.applyBecomeMediator(stub, args)
	case "HandleForApplyBecomeMediator":
		//基金会对加入申请Mediator进行处理
		return d.handleForApplyBecomeMediator(stub, args)
	case "MediatorApplyQuitMediator":
		//申请退出Mediator
		return d.mediatorApplyQuitMediator(stub, args)
	case "HandleForApplyQuitMediator":
		//基金会对退出申请Mediator进行处理
		return d.handleForApplyQuitMediator(stub, args)
	case "MediatorPayToDepositContract":
		//mediator 交付保证金
		return d.mediatorPayToDepositContract(stub, args)
	case "JuryPayToDepositContract":
		//jury 交付保证金
		return d.juryPayToDepositContract(stub, args)
	case "DeveloperPayToDepositContract":
		//developer 交付保证金
		return d.developerPayToDepositContract(stub, args)
	case "MediatorApplyCashback":
		//mediator 申请提取保证金
		return d.mediatorApplyCashback(stub, args)
	case "HandleForMediatorApplyCashback":
		//基金会处理提取保证金
		return d.handleForMediatorApplyCashback(stub, args)
	case "JuryApplyCashback":
		//jury 申请提取保证金
		return d.juryApplyCashback(stub, args)
	case "HandleForJuryApplyCashback":
		//基金会处理提取保证金
		return d.handleForJuryApplyCashback(stub, args)
	case "DeveloperApplyCashback":
		//developer 申请提取保证金
		return d.developerApplyCashback(stub, args)
	case "HandleForDeveloperApplyCashback":
		//基金会处理提取保证金
		return d.handleForDeveloperApplyCashback(stub, args)
	//case "ApplyForDepositCashback":
	//	//申请保证金退还
	//	return d.applyForDepositCashback(stub, args)
	//case "HandleForCashbackApplication":
	//	//基金会对申请做相应的处理
	//	return d.handleApplications(stub, args, "Cashback")
	case "ApplyForForfeitureDeposit":
		//申请保证金没收
		//void forfeiture_deposit(const witness_object& wit, token_type amount)
		return d.applyForForfeitureDeposit(stub, args)
	case "HandleForForfeitureApplication":
		//基金会对申请做相应的处理
		return d.handleForForfeitureApplication(stub, args)
	//获取提取保证金申请列表
	case "GetListForCashbackApplication":
		list, err := stub.GetState("ListForCashback")
		if err != nil {
			return shim.Error(err.Error())
		}
		if list == nil {
			return shim.Success([]byte("[]"))
		}
		return shim.Success(list)
		//获取没收保证金申请列表
	case "GetListForForfeitureApplication":
		list, err := stub.GetState("ListForForfeiture")
		if err != nil {
			return shim.Error(err.Error())
		}
		if list == nil {
			return shim.Success([]byte("[]"))
		}
		return shim.Success(list)
		//获取Mediator候选列表
	case "GetListForMediatorCandidate":
		list, err := stub.GetState("MediatorList")
		if err != nil {
			return shim.Error(err.Error())
		}
		if list == nil {
			return shim.Success([]byte("[]"))
		}
		return shim.Success(list)
		//获取Jury候选列表
	case "GetListForJuryCandidater":
		list, err := stub.GetState("JuryList")
		if err != nil {
			return shim.Error(err.Error())
		}
		if list == nil {
			return shim.Success([]byte("[]"))
		}
		return shim.Success(list)
		//获取Contract Developer候选列表
	case "GetListForDeveloperCandidate":
		list, err := stub.GetState("DeveloperList")
		if err != nil {
			return shim.Error(err.Error())
		}
		if list == nil {
			return shim.Success([]byte("[]"))
		}
		return shim.Success(list)
		//获取某个节点的账户
	case "GetCandidateBalanceWithAddr":
		list, err := stub.GetState(args[0])
		if err != nil {
			return shim.Error(err.Error())
		}
		if list == nil {
			return shim.Success([]byte("[]"))
		}
		return shim.Success(list)
		//获取Mediator申请加入列表
	case "GetBecomeMediatorApplyList":
		list, err := stub.GetState("ListForApplyBecomeMediator")
		if err != nil {
			return shim.Error(err.Error())
		}
		if list == nil {
			return shim.Success([]byte("[]"))
		}
		return shim.Success(list)
		//获取已同意的mediator列表
	case "GetAgreeForBecomeMediatorList":
		list, err := stub.GetState("ListForAgreeBecomeMediator")
		if err != nil {
			return shim.Error(err.Error())
		}
		if list == nil {
			return shim.Success([]byte("[]"))
		}
		return shim.Success(list)
		//获取Mediator申请退出列表
	case "GetQuitMediatorApplyList":
		list, err := stub.GetState("ListForApplyQuitMediator")
		if err != nil {
			return shim.Error(err.Error())
		}
		if list == nil {
			return shim.Success([]byte("[]"))
		}
		return shim.Success(list)
	}
	return shim.Success([]byte("Invoke error"))
}

//申请加入Mediator
func (d *DepositChaincode) applyBecomeMediator(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return applyBecomeMediator(stub, args)
}

//基金会对申请加入Mediator进行处理
func (d *DepositChaincode) handleForApplyBecomeMediator(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return handleForApplyBecomeMediator(stub, args)
}

//申请退出Mediator
func (d *DepositChaincode) mediatorApplyQuitMediator(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return mediatorApplyQuitMediator(stub, args)
}

//基金会对申请退出Mediator进行处理
func (d *DepositChaincode) handleForApplyQuitMediator(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return handleForApplyQuitMediator(stub, args)
}

//交付保证金
//handle witness pay
//func (d *DepositChaincode) depositWitnessPay(stub shim.ChaincodeStubInterface, args []string) pb.Response {
//	//第一个参数：合约地址；第二个参数：保证金；第三个参数：角色（Mediator Jury ContractDeveloper)
//	//Deposit("contractAddr","2000","Mediator")
//	if len(args) != 1 {
//		return shim.Error("Input parameter Success,need one parameter.")
//	}
//	//获取 请求 调用 地址（即交付保证节点地址）
//	invokeAddr, err := stub.GetInvokeAddress()
//	if err != nil {
//		return shim.Error("GetInvokeFromAddr error:")
//	}
//	fmt.Println("invokeFromAddr address = ", invokeAddr)
//	//获取 请求 ptn 数量（即交付保证金数量）
//	invokeTokens, err := stub.GetInvokeTokens()
//	if err != nil {
//		return shim.Success([]byte("GetPayToContractPtnTokens error:"))
//	}
//	fmt.Printf("invokeTokens=%v", invokeTokens)
//	//获取退保证金数量，将 string 转 uint64
//	//TODO test
//	//ptnAccount, _ := strconv.ParseUint(args[0], 10, 64)
//	//invokeTokens.Amount = ptnAccount
//	//fmt.Println("invokeTokens ", invokeTokens.Amount)
//	//fmt.Printf("invokeTokens %#v\n", invokeTokens.Asset)
//	//获取角色
//	role := args[0]
//	switch {
//	//case role == "Mediator":
//	//	//处理Mediator交付保证金
//	//	return d.handleMediatorDepositWitnessPay(stub, invokeAddr, invokeTokens)
//	case role == "Jury":
//		//处理Jury交付保证金
//		return d.handleJuryDepositWitnessPay(stub, invokeAddr, invokeTokens)
//	case role == "Developer":
//		//处理Developer交付保证金
//		return d.handleDeveloperDepositWitnessPay(stub, invokeAddr, invokeTokens)
//	default:
//		return shim.Success([]byte("role error."))
//	}
//}

//处理 Mediator
//func (d *DepositChaincode) handleMediatorDepositWitnessPay(stub shim.ChaincodeStubInterface, invokeAddr string, invokeTokens *modules.InvokeTokens) pb.Response {
//	//获取同意列表
//	agreeList, err := stub.GetAgreeForBecomeMediatorList()
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	if agreeList == nil {
//		return shim.Error("Your node does not in agree list for mediator.")
//	}
//	//获取节点信息
//	mediator := &modules.MediatorInfo{}
//	for _, m := range agreeList {
//		if strings.Compare(m.Address, invokeAddr) == 0 {
//			mediator = m
//		}
//	}
//	if mediator == nil {
//		return shim.Error("Your node does not in agree list for mediator.")
//	}
//	mediator.Time = time.Now().UTC()
//	//获取一下该用户下的账簿情况
//	balance, err := stub.GetDepositBalance(invokeAddr)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//
//	//stateValueBytes, err := stub.GetState(invokeAddr)
//	//if err != nil {
//	//	return shim.Success([]byte("Get account balance from ledger error:"))
//	//}
//	//stateValues := new(modules.DepositStateValues)
//	//
//	//stateValue := new(modules.DepositStateValue)
//	//账户不存在，第一次参与
//	if balance == nil {
//		//判断保证金是否足够(Mediator第一次交付必须足够)
//		if invokeTokens.Amount < depositAmountsForMediator {
//			//TODO 第一次交付不够的话，这里必须终止
//			return shim.Error("Payment amount is insufficient.")
//		}
//		//加入列表
//		//addList("Mediator", invokeAddr, stub)
//		err = addCandidateListForMediator(stub, mediator)
//		if err != nil {
//			return shim.Error(err.Error())
//		}
//		balance = new(modules.DepositBalance)
//		//处理数据
//		balance.EnterTime = time.Now().UTC()
//		d.updateForPayValue(balance, invokeTokens)
//	} else {
//		//已经是mediator了
//		//err = json.Unmarshal(stateValueBytes, stateValues)
//		//if err != nil {
//		//	return shim.Success([]byte("Unmarshal stateValueBytes error"))
//		//}
//		//TODO 再次交付保证金时，先计算当前余额的币龄奖励
//		awards := award.GetAwardsWithCoins(balance.TotalAmount, balance.LastModifyTime.Unix())
//		balance.TotalAmount += awards
//		//获取币龄
//		//endTime := time.Now().UTC()
//		//coinDays := award.GetCoinDay(stateValues.TotalAmount, stateValues.LastModifyTime, endTime)
//		////计算币龄收益
//		//awards := award.CalculateAwardsForDepositContractNodes(coinDays)
//
//		//处理数据
//		d.updateForPayValue(balance, invokeTokens)
//	}
//	//对结果序列化并更新数据
//	return d.marshalForBalance(stub, invokeAddr, balance)
//}

//处理 Jury
//func (d *DepositChaincode) handleJuryDepositWitnessPay(stub shim.ChaincodeStubInterface, invokeAddr string, invokeTokens *modules.InvokeTokens) pb.Response {
//	//获取一下该用户下的账簿情况
//	balance, err := stub.GetDepositBalance(invokeAddr)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	isJury := false
//	if balance == nil {
//		balance = new(modules.DepositBalance)
//		if invokeTokens.Amount >= depositAmountsForJury {
//			//加入列表
//			//addList("Jury", invokeAddr, stub)
//			err = addCandaditeList(invokeAddr, stub, "JuryList")
//			if err != nil {
//				return shim.Error(err.Error())
//			}
//			isJury = true
//			balance.EnterTime = time.Now().UTC()
//		}
//		//处理数据
//		//fmt.Println(balance)
//		//fmt.Println(invokeTokens)
//		//return shim.Success([]byte("ok"))
//		d.updateForPayValue(balance, invokeTokens)
//	} else {
//		//账户已存在，进行信息的更新操作
//		if balance.TotalAmount >= depositAmountsForJury {
//			//原来就是jury
//			isJury = true
//			//TODO 再次交付保证金时，先计算当前余额的币龄奖励
//			awards := award.GetAwardsWithCoins(balance.TotalAmount, balance.LastModifyTime.Unix())
//			balance.TotalAmount += awards
//
//		}
//		//处理交付保证金数据
//		d.updateForPayValue(balance, invokeTokens)
//	}
//	if !isJury {
//		//判断交了保证金后是否超过了jury
//		if balance.TotalAmount >= depositAmountsForJury {
//			//addList("Jury", invokeAddr, stub)
//			err = addCandaditeList(invokeAddr, stub, "JuryList")
//			if err != nil {
//				return shim.Error(err.Error())
//			}
//			balance.EnterTime = time.Now().UTC()
//		}
//	}
//	//对结果序列化并更新数据
//	return d.marshalForBalance(stub, invokeAddr, balance)
//}

//处理 ContractDeveloper
//func (d *DepositChaincode) handleDeveloperDepositWitnessPay(stub shim.ChaincodeStubInterface, invokeAddr string, invokeTokens *modules.InvokeTokens) pb.Response {
//	//获取一下该用户下的账簿情况
//	balance, err := stub.GetDepositBalance(invokeAddr)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	isDeveloper := false
//	if balance == nil {
//		balance = new(modules.DepositBalance)
//		if invokeTokens.Amount >= depositAmountsForDeveloper {
//			//加入列表
//			//addList("Developer", invokeAddr, stub)
//			err = addCandaditeList(invokeAddr, stub, "DeveloperList")
//			if err != nil {
//				return shim.Error(err.Error())
//			}
//			isDeveloper = true
//			balance.EnterTime = time.Now().UTC()
//		}
//		//处理数据
//		d.updateForPayValue(balance, invokeTokens)
//	} else {
//		//账户已存在，进行信息的更新操作
//		if balance.TotalAmount >= depositAmountsForDeveloper {
//			//原来就是jury
//			isDeveloper = true
//			//TODO 再次交付保证金时，先计算当前余额的币龄奖励
//			awards := award.GetAwardsWithCoins(balance.TotalAmount, balance.LastModifyTime.Unix())
//			balance.TotalAmount += awards
//
//		}
//		//处理交付保证金数据
//		d.updateForPayValue(balance, invokeTokens)
//	}
//	if !isDeveloper {
//		//判断交了保证金后是否超过了jury
//		if balance.TotalAmount >= depositAmountsForDeveloper {
//			//addList("Developer", invokeAddr, stub)
//			err = addCandaditeList(invokeAddr, stub, "DeveloperList")
//			if err != nil {
//				return shim.Error(err.Error())
//			}
//			balance.EnterTime = time.Now().UTC()
//		}
//	}
//	//对结果序列化并更新数据
//	return d.marshalForBalance(stub, invokeAddr, balance)
//}

//处理交付保证金数据
//func (d *DepositChaincode) updateForPayValue(balance *modules.DepositBalance, invokeTokens *modules.InvokeTokens) {
//	balance.TotalAmount += invokeTokens.Amount
//	balance.LastModifyTime = time.Now().UTC()
//
//	payTokens := &modules.InvokeTokens{}
//	payValue := &modules.PayValue{PayTokens: payTokens}
//	payValue.PayTokens.Amount = invokeTokens.Amount
//	payValue.PayTokens.Asset = invokeTokens.Asset
//	payValue.PayTime = time.Now().UTC()
//
//	balance.PayValues = append(balance.PayValues, payValue)
//}

//对结果序列化并更新数据
func (d *DepositChaincode) marshalForBalance(stub shim.ChaincodeStubInterface, nodeAddr string, balance *modules.DepositBalance) pb.Response {
	balanceByte, err := json.Marshal(balance)
	if err != nil {
		return shim.Success([]byte("Marshal balance error" + err.Error()))
	}
	err = stub.PutState(nodeAddr, balanceByte)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("ok"))
}

//保证金退还，只申请，当然符合要求了才能申请成功，并且加入申请列表
//handle cashback rewards
//func (d *DepositChaincode) applyForDepositCashback(stub shim.ChaincodeStubInterface, args []string) pb.Response {
//	//第一个参数：数量；第二个参数：角色（角色（Mediator Jury ContractDeveloper)
//	//depositCashback("保证金数量","Mediator")
//	if len(args) != 2 {
//		return shim.Success([]byte("Input parameter Success,need two parameters."))
//	}
//	//获取 请求 调用 地址
//	invokeAddr, err := stub.GetInvokeAddress()
//	if err != nil {
//		return shim.Success([]byte("GetInvokeFromAddr error:"))
//	}
//	fmt.Println("invokeAddr address ", invokeAddr)
//	//获取退保证金数量，将 string 转 uint64
//	ptnAccount, err := strconv.ParseUint(args[0], 10, 64)
//	if err != nil {
//		return shim.Success([]byte("String transform to uint64 error:"))
//	}
//	fmt.Println("ptnAccount  args[0] ", ptnAccount)
//	asset := modules.NewPTNAsset()
//	invokeTokens := &modules.InvokeTokens{
//		Amount: ptnAccount,
//		Asset:  asset,
//	}
//	//
//	////TODO 先获取数据库信息
//	balance, err := stub.GetDepositBalance(invokeAddr)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	if balance == nil {
//		return shim.Error("账户不存在")
//	}
//	//stateValueBytes, err := stub.GetState(invokeAddr)
//	//if err != nil {
//	//	return shim.Success([]byte("Get account balance from ledger error:"))
//	//}
//	////判断数据库是否为空
//	//if stateValueBytes == nil {
//	//	return shim.Success([]byte("Your account does not exist."))
//	//}
//	//balanceValue := new(modules.DepositStateValues)
//	////如果不为空，反序列化数据库信息
//	//err = json.Unmarshal(stateValueBytes, balanceValue)
//	//if err != nil {
//	//	return shim.Success([]byte("Unmarshal stateValueBytes error:"))
//	//}
//	//比较退款数量和数据库数量
//	//Asset判断
//	//数量比较
//	if balance.TotalAmount < invokeTokens.Amount {
//		return shim.Success([]byte("Your delivery amount with ptn token is insufficient."))
//	}
//	err = d.addListForCashback(args[1], stub, invokeAddr, invokeTokens)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	return shim.Success([]byte("申请成功"))
//}

//加入退款申请列表
//func (d *DepositChaincode) addListForCashback(role string, stub shim.ChaincodeStubInterface, invokeAddr string, invokeTokens *modules.InvokeTokens) error {
//	//先获取申请列表
//	listForCashback, err := stub.GetListForCashback()
//	if err != nil {
//		return err
//	}
//	////序列化
//	cashback := new(modules.Cashback)
//	cashback.CashbackAddress = invokeAddr
//	cashback.CashbackTokens = invokeTokens
//	cashback.Role = role
//	cashback.CashbackTime = time.Now().UTC().Unix()
//	if listForCashback == nil {
//		listForCashback = []*modules.Cashback{}
//		listForCashback = append(listForCashback, cashback)
//	} else {
//		listForCashback = append(listForCashback, cashback)
//	}
//	//反序列化
//	listForCashbackByte, err := json.Marshal(listForCashback)
//	if err != nil {
//		return err
//	}
//	err = stub.PutState("ListForCashback", listForCashbackByte)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (d *DepositChaincode) handleForForfeitureApplication(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//地址，申请时间，是否同意
	if len(args) != 3 {
		return shim.Error("Input parameter error,need three parameters.")
	}

	//基金会地址
	invokeAddr, _ := stub.GetInvokeAddress()
	fmt.Println("invokeAddr==", invokeAddr)
	//获取系统配置基金会地址
	//foundationAddress, err := stub.GetSystemConfig("FoundationAddress")
	//if err != nil {
	//	return shim.Success([]byte("获取基金会地址错误"))
	//}
	//fmt.Println("foundationAddress==", foundationAddress)
	//判断没收请求地址是否是基金会地址
	if strings.Compare(invokeAddr, foundationAddress) != 0 {
		return shim.Error("请求地址不正确，请使用基金会的地址")
	}

	//获取没收节点地址
	//nodeAddr, err := common.StringToAddress(args[0])
	//if err != nil {
	//	return shim.Success([]byte("string to address error"))
	//}
	fmt.Println("nodeAddr ", args[0])
	//获取一下该用户下的账簿情况
	addr := args[0]
	balance, err := stub.GetDepositBalance(addr)
	if err != nil {
		return shim.Error(err.Error())
	}
	//判断没收节点账户是否为空
	if balance == nil {
		return shim.Error("you have not depositWitnessPay for deposit.")
	}
	////获取没收节点的账本信息
	//stateValueBytes, err := stub.GetState(nodeAddr)
	//if err != nil {
	//	return shim.Success([]byte("Get account balance from ledger error:"))
	//}
	//
	//balanceValue := new(modules.DepositStateValues)
	////将没收节点账户序列化
	//err = json.Unmarshal(stateValueBytes, balanceValue)
	//if err != nil {
	//	return shim.Success([]byte("unmarshal accBalByte error"))
	//}

	//获取申请时间戳
	applyTime, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return shim.Error("string to int64 error " + err.Error())
	}
	fmt.Println("applytime ", applyTime)
	//获取是否同意
	check := args[2]

	//获取处理类别
	//switch {
	//case applyTye == "Cashback":
	//	return d.handleDepositCashbackApplication(stub, invokeAddr, args[0], applyTime, balance, check)
	//case applyTye == "Forfeiture":
	return d.handleForfeitureDepositApplication(stub, invokeAddr, addr, applyTime, balance, check)
	//default:
	//	return shim.Error("类别错误")
	//}
}

//这里是基金会处理保证金提取的请求
//func (d *DepositChaincode) handleDepositCashbackApplication(stub shim.ChaincodeStubInterface, foundationAddr, cashbackAddr string, applyTime int64, balance *modules.DepositBalance, check string) pb.Response {
//	//提取保证金节点地址，申请时间
//	if check == "ok" {
//		return d.agreeForApplyCashback(stub, foundationAddr, cashbackAddr, applyTime, balance)
//	} else if check == "no" {
//		return d.disagreeForApplyCashback(stub, cashbackAddr, applyTime)
//	}
//	return shim.Success([]byte("ok"))
//}

//同意申请退保证金请求
//func (d *DepositChaincode) agreeForApplyCashback(stub shim.ChaincodeStubInterface, foundationAddr, cashbackAddr string, applyTime int64, balance *modules.DepositBalance) pb.Response {
//	//获取请求列表
//	listForCashback, err := stub.GetListForCashback()
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	if listForCashback == nil {
//		return shim.Error("listForCashback is nil.")
//	}
//	//在申请退款保证金列表中移除该节点
//	//fmt.Println(listForCashback)
//	//fmt.Println(cashbackAddr)
//	//fmt.Println(applyTime)
//	cashbackValue, err := moveInApplyForCashbackList(stub, listForCashback, cashbackAddr, applyTime)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	//fmt.Println(cashbackValue)
//	//fmt.Printf("%#v\n\n", cashbackValue)
//	if cashbackValue == nil {
//		return shim.Error("列表里没有该申请")
//	}
//	//还得判断一下是否超过余额
//	if cashbackValue.CashbackTokens.Amount > balance.TotalAmount {
//		return shim.Error("退款大于账户余额")
//	}
//	role := cashbackValue.Role
//	//判断节点类型
//	switch {
//	case role == "Mediator":
//		return d.handleMediatorDepositCashback(foundationAddr, cashbackAddr, cashbackValue, balance, stub)
//	case role == "Jury":
//		return d.handleJuryDepositCashback(stub, cashbackAddr, cashbackValue, balance)
//	case role == "Developer":
//		return d.handleDeveloperDepositCashback(stub, cashbackAddr, cashbackValue, balance)
//	default:
//		return shim.Error("role error")
//	}
//}

//退保证金请求
//func (d *DepositChaincode) handleMediatorDepositCashback(foundationAddr, cashbackAddr string, cashbackValue *modules.Cashback, balance *modules.DepositBalance, stub shim.ChaincodeStubInterface) pb.Response {
//	var err error
//	//规定mediator 退款要么全部退，要么退款后，剩余数量在mediator数量范围内，
//	//计算余额
//	result := balance.TotalAmount - cashbackValue.CashbackTokens.Amount
//	//判断是否全部退
//	if result == 0 {
//		//加入候选列表的时的时间
//		startTime := balance.EnterTime.YearDay()
//		//当前时间
//		endTime := time.Now().UTC().YearDay()
//		//判断是否已超过规定周期
//		if endTime-startTime >= depositPeriod {
//			//退出全部，即删除cashback
//			err = d.cashbackAllDeposit("MediatorList", stub, cashbackAddr, cashbackValue.CashbackTokens, balance)
//			if err != nil {
//				return shim.Success([]byte(err.Error()))
//			}
//			return shim.Success([]byte("成功退出"))
//		} else {
//			//没有超过周期，不能退出
//			return shim.Error("还在规定周期之内，不得退出列表")
//		}
//	} else if result < depositAmountsForMediator {
//		//说明退款后，余额少于规定数量
//		return shim.Error("说明退款后，余额少于规定数量，对于Mediator来说，如果退部分保证后，余额少于规定数量，则不允许提款或者没收")
//	} else {
//		//TODO 这是只退一部分钱，剩下余额还是在规定范围之内
//		return d.cashbackSomeDeposit("Mediator", stub, cashbackAddr, cashbackValue, balance)
//	}
//}
//
////对Jury退保证金的处理
//func (d *DepositChaincode) handleJuryDepositCashback(stub shim.ChaincodeStubInterface, cashbackAddr string, cashbackValue *modules.Cashback, balance *modules.DepositBalance) pb.Response {
//	var res pb.Response
//	if balance.TotalAmount >= depositAmountsForJury {
//		//已在列表中
//		res = d.handleJuryFromList(stub, cashbackAddr, cashbackValue, balance)
//	} else {
//		////TODO 不在列表中,没有奖励，直接退
//		res = d.handleCommonJuryOrDev(stub, cashbackAddr, cashbackValue, balance)
//	}
//	return res
//}

////Jury已在列表中
//func (d *DepositChaincode) handleJuryFromList(stub shim.ChaincodeStubInterface, cashbackAddr string, cashbackValue *modules.Cashback, balance *modules.DepositBalance) pb.Response {
//	//退出列表
//	var err error
//	//计算余额
//	resule := balance.TotalAmount - cashbackValue.CashbackTokens.Amount
//	//判断是否退出列表
//	if resule == 0 {
//		//加入列表时的时间
//		startTime := balance.EnterTime.YearDay()
//		//当前退出时间
//		endTime := time.Now().UTC().YearDay()
//		//判断是否已到期
//		if endTime-startTime >= depositPeriod {
//			//退出全部，即删除cashback，利息计算好了
//			err = d.cashbackAllDeposit("Jury", stub, cashbackAddr, cashbackValue.CashbackTokens, balance)
//			if err != nil {
//				return shim.Success([]byte(err.Error()))
//			}
//			return shim.Success([]byte("成功退出"))
//		} else {
//			return shim.Success([]byte("未到期，不能退出列表"))
//		}
//	} else {
//		//TODO 退出一部分，且退出该部分金额后还在列表中，还没有计算利息
//		//d.addListForCashback("Jury", stub, cashbackAddr, invokeTokens)
//		return d.cashbackSomeDeposit("Jury", stub, cashbackAddr, cashbackValue, balance)
//	}
//}
//
////对Developer退保证金的处理
//func (d *DepositChaincode) handleDeveloperDepositCashback(stub shim.ChaincodeStubInterface, cashbackAddr string, cashbackValue *modules.Cashback, balance *modules.DepositBalance) pb.Response {
//	var res pb.Response
//	if balance.TotalAmount >= depositAmountsForDeveloper {
//		//已在列表中
//		res = d.handleDeveloperFromList(stub, cashbackAddr, cashbackValue, balance)
//	} else {
//		////TODO 不在列表中,没有奖励，直接退
//		res = d.handleCommonJuryOrDev(stub, cashbackAddr, cashbackValue, balance)
//	}
//	return res
//}
//
////Jury or developer 可以随时退保证金，只是不在列表中的话，没有奖励
//func (d *DepositChaincode) handleCommonJuryOrDev(stub shim.ChaincodeStubInterface, cashbackAddr string, cashbackValue *modules.Cashback, balance *modules.DepositBalance) pb.Response {
//	//调用从合约把token转到请求地址
//	err := stub.PayOutToken(cashbackAddr, cashbackValue.CashbackTokens, 0)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	//fmt.Printf("balanceValue=%s\n", balanceValue)
//	//v := handleValues(balanceValue.Values, tokens)
//	//balanceValue.Values = v
//	balance.LastModifyTime = time.Now().UTC()
//	balance.TotalAmount -= cashbackValue.CashbackTokens.Amount
//	//fmt.Printf("balanceValue=%s\n", balanceValue)
//	//TODO
//	balance.CashbackValues = append(balance.CashbackValues, cashbackValue)
//
//	return d.marshalForBalance(stub, cashbackAddr, balance)
//}
//
////Developer已在列表中
//func (d *DepositChaincode) handleDeveloperFromList(stub shim.ChaincodeStubInterface, cashbackAddr string, cashbackValue *modules.Cashback, balance *modules.DepositBalance) pb.Response {
//	//退出列表
//	var err error
//	//计算余额
//	result := balance.TotalAmount - cashbackValue.CashbackTokens.Amount
//	//判断是否退出列表
//	if result == 0 {
//		//加入列表时的时间
//		startTime := balance.EnterTime.YearDay()
//		//当前退出时间
//		endTime := time.Now().UTC().YearDay()
//		//判断是否已到期
//		if endTime-startTime >= depositPeriod {
//			//退出全部，即删除cashback，利息计算好了
//			err = d.cashbackAllDeposit("Developer", stub, cashbackAddr, cashbackValue.CashbackTokens, balance)
//			if err != nil {
//				return shim.Success([]byte(err.Error()))
//			}
//			return shim.Success([]byte("成功退出"))
//		} else {
//			return shim.Success([]byte("未到期，不能退出列表"))
//		}
//	} else {
//		//TODO 退出一部分，且退出该部分金额后还在列表中，还没有计算利息
//		//d.addListForCashback("Jury", stub, cashbackAddr, invokeTokens)
//		return d.cashbackSomeDeposit("Developer", stub, cashbackAddr, cashbackValue, balance)
//	}
//}

//处理申请提保证金请求并移除列表
//func (d *DepositChaincode) cashbackAllDeposit(role string, stub shim.ChaincodeStubInterface, cashbackAddr string, invokeTokens *modules.InvokeTokens, balance *modules.DepositBalance) error {
//	//计算保证金全部利息
//	//获取币龄
//	endTime := time.Now().UTC()
//	coinDays := award.GetCoinDay(balance.TotalAmount, balance.LastModifyTime, endTime)
//	//计算币龄收益
//	awards := award.CalculateAwardsForDepositContractNodes(coinDays)
//	//本金+利息
//	invokeTokens.Amount += awards
//	//调用从合约把token转到请求地址
//	err := stub.PayOutToken(cashbackAddr, invokeTokens, 0)
//	if err != nil {
//		return err
//	}
//	//移除出列表
//	err = moveCandidate(role, cashbackAddr, stub)
//	if err != nil {
//		return err
//	}
//	//删除节点
//	err = stub.DelState(cashbackAddr)
//	if err != nil {
//		return err
//	}
//	return nil
//}

//不需要移除候选列表，但是要没收一部分保证金
//func (d *DepositChaincode) cashbackSomeDeposit(role string, stub shim.ChaincodeStubInterface, cashbackAddr string, cashbackValue *modules.Cashback, balance *modules.DepositBalance) pb.Response {
//	//tokens.Amount += awards
//	//调用从合约把token转到请求地址
//	err := stub.PayOutToken(cashbackAddr, cashbackValue.CashbackTokens, 0)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	awards := award.GetAwardsWithCoins(balance.TotalAmount, balance.LastModifyTime.Unix())
//	fmt.Println("awards ", awards)
//	balance.LastModifyTime = time.Now().UTC()
//	//加上利息奖励
//	balance.TotalAmount += awards
//	//减去提取部分
//	balance.TotalAmount -= cashbackValue.CashbackTokens.Amount
//	//TODO 如果推出后低于保证金，则退出列表
//	if role == "Jury" {
//		//如果推出后低于保证金，则退出列表
//		if balance.TotalAmount < depositAmountsForJury {
//			//handleMember("Jury", cashbackAddr, stub)
//			err = moveCandidate("JuryList", cashbackAddr, stub)
//			if err != nil {
//				return shim.Error(err.Error())
//			}
//		}
//	} else if role == "Developer" {
//		//如果推出后低于保证金，则退出列表
//		if balance.TotalAmount < depositAmountsForDeveloper {
//			//handleMember("Developer", cashbackAddr, stub)
//			err = moveCandidate("DeveloperList", cashbackAddr, stub)
//			if err != nil {
//				return shim.Error(err.Error())
//			}
//		}
//	}
//	//TODO 加入提款记录
//	balance.CashbackValues = append(balance.CashbackValues, cashbackValue)
//	//序列化
//	return d.marshalForBalance(stub, cashbackAddr, balance)
//}

//社区申请没收某节点的保证金数量
func (d DepositChaincode) applyForForfeitureDeposit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//没收地址 数量 角色 额外说明
	//forfeiture string, invokeTokens modules.InvokeTokens, role, extra string
	if len(args) != 3 {
		return shim.Error("need three parameters.")
	}
	//申请地址
	invokeAddr, _ := stub.GetInvokeAddress()
	forfeiture := new(modules.Forfeiture)
	forfeiture.ApplyAddress = invokeAddr
	//forfeitureAddr, err := common.StringToAddress(args[0])
	////获取没收节点地址
	//if err != nil {
	//	return shim.Success([]byte("string to address error"))
	//}
	//fmt.Println(forfeitureAddr)
	//TODO 获取没收节点的账本信息
	//stateValueBytes, err := stub.GetState(forfeitureAddr)
	//if err != nil {
	//	return shim.Success([]byte("Get account balance from ledger error:"))
	//}
	////判断没收节点账户是否为空
	//if stateValueBytes == nil {
	//	return shim.Success([]byte("you have not depositWitnessPay for deposit."))
	//}
	//balanceValue := new(modules.DepositStateValues)
	////将没收节点账户序列化
	//err = json.Unmarshal(stateValueBytes, balanceValue)
	//if err != nil {
	//	return shim.Success([]byte("unmarshal accBalByte error"))
	//}
	//获取没收保证金数量，将 string 转 uint64
	ptnAccount, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		return shim.Success([]byte("String transform to uint64 error:"))
	}
	fmt.Println("ptnAccount  args[1] ", ptnAccount)
	//判断账户余额和没收请求数量
	//if balanceValue.TotalAmount < ptnAccount {
	//	return shim.Success([]byte("Forfeiture too many."))
	//}
	forfeiture.ForfeitureAddress = args[0]
	asset := modules.NewPTNAsset()
	invokeTokens := &modules.InvokeTokens{
		Amount: ptnAccount,
		Asset:  asset,
	}
	forfeiture.ApplyTokens = invokeTokens
	forfeiture.ForfeitureRole = args[2]
	forfeiture.ApplyTime = time.Now().UTC().Unix()
	//先获取列表，再更新列表
	listForForfeiture, err := stub.GetListForForfeiture()
	if err != nil {
		return shim.Error(err.Error())
	}
	if listForForfeiture == nil {
		listForForfeiture = []*modules.Forfeiture{forfeiture}
	} else {
		isExist := isInForfeiturelist(forfeiture.ForfeitureAddress, listForForfeiture)
		if isExist {
			return shim.Error("node is exist in the list.")
		}
		listForForfeiture = append(listForForfeiture, forfeiture)
	}
	listForForfeitureByte, err := json.Marshal(listForForfeiture)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState("ListForForfeiture", listForForfeitureByte)
	return shim.Success([]byte("申请成功"))
}

//查找节点是否在列表中
func isInForfeiturelist(addr string, list []*modules.Forfeiture) bool {
	for _, m := range list {
		if strings.Compare(addr, m.ForfeitureAddress) == 0 {
			return true
		}
	}
	return false
}

//基金会处理没收请求
func (d *DepositChaincode) handleForfeitureDepositApplication(stub shim.ChaincodeStubInterface, foundationAddr, forfeitureAddr string, applyTime int64, balance *modules.DepositBalance, check string) pb.Response {
	//check 如果为ok，则同意此申请，如果为no，则不同意此申请
	if check == "ok" {
		return d.agreeForApplyForfeiture(stub, foundationAddr, forfeitureAddr, applyTime, balance)
	} else if check == "no" {
		//移除申请列表，不做处理
		return d.disagreeForApplyForfeiture(stub, forfeitureAddr, applyTime)
	}
	return shim.Error("请确认是否同意")
}

//不同意提取请求，则直接从提保证金列表中移除该节点
func (d *DepositChaincode) disagreeForApplyCashback(stub shim.ChaincodeStubInterface, cashbackAddr string, applyTime int64) pb.Response {
	//获取没收列表
	listForCashback, err := stub.GetListForCashback()
	if err != nil {
		return shim.Error(err.Error())
	}
	if listForCashback == nil {
		return shim.Error("listForCashback is nil")
	}
	fmt.Println("moveInApplyForCashbackList==>", listForCashback)
	node, err := moveInApplyForCashbackList(stub, listForCashback, cashbackAddr, applyTime)
	if err != nil {
		return shim.Error(err.Error())
	}
	if node == nil {
		return shim.Error("列表里没有该申请")
	}
	fmt.Println("moveInApplyForCashbackList==>", listForCashback)
	return shim.Success([]byte("移除列表成功"))
}

//不同意这样没收请求，则直接从没收列表中移除该节点
func (d *DepositChaincode) disagreeForApplyForfeiture(stub shim.ChaincodeStubInterface, forfeiture string, applyTime int64) pb.Response {
	//获取没收列表
	listForForfeiture, err := stub.GetListForForfeiture()
	if err != nil {
		return shim.Error(err.Error())
	}
	if listForForfeiture == nil {
		return shim.Error("listForForfeiture is nil.")
	}
	isExist := isInForfeiturelist(forfeiture, listForForfeiture)
	if !isExist {
		return shim.Error("node is not exist in the list.")
	}
	_, err = moveInApplyForForfeitureList(stub, listForForfeiture, forfeiture, applyTime)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("移除列表成功"))
}

//同意申请没收请求
func (d *DepositChaincode) agreeForApplyForfeiture(stub shim.ChaincodeStubInterface, foundationAddr, forfeitureAddr string, applyTime int64, balance *modules.DepositBalance) pb.Response {
	//获取列表
	listForForfeiture, err := stub.GetListForForfeiture()
	if err != nil {
		return shim.Error(err.Error())
	}
	if listForForfeiture == nil {
		return shim.Error("listForForfeiture is nil.")
	}
	isExist := isInForfeiturelist(forfeitureAddr, listForForfeiture)
	if !isExist {
		return shim.Error("node is not exist in the list.")
	}
	//在列表中移除，并获取没收情况
	forfeiture, err := moveInApplyForForfeitureList(stub, listForForfeiture, forfeitureAddr, applyTime)
	if err != nil {
		return shim.Error(err.Error())
	}
	//判断余额
	if forfeiture.ApplyTokens.Amount > balance.TotalAmount {
		return shim.Error("没收数量超过余额")
	}
	//判断节点类型
	switch {
	case forfeiture.ForfeitureRole == "Mediator":
		return d.handleMediatorForfeitureDeposit(foundationAddr, forfeiture, balance, stub)
	case forfeiture.ForfeitureRole == "Jury":
		return d.handleJuryForfeitureDeposit(foundationAddr, forfeiture, balance, stub)
	case forfeiture.ForfeitureRole == "Developer":
		return d.handleDeveloperForfeitureDeposit(foundationAddr, forfeiture, balance, stub)
	default:
		return shim.Error("role error")
	}
}

//处理申请没收请求并移除列表
func (d *DepositChaincode) forfeitureAllDeposit(role string, stub shim.ChaincodeStubInterface, foundationAddr, forfeitureAddr string, invokeTokens *modules.InvokeTokens) error {
	//TODO 没收保证金是否需要计算利息
	//调用从合约把token转到请求地址
	err := stub.PayOutToken(foundationAddr, invokeTokens, 0)
	if err != nil {
		return err
	}
	//移除出列表
	//handleMember(role, forfeitureAddr, stub)
	err = moveCandidate(role, forfeitureAddr, stub)
	if err != nil {
		return err
	}
	//删除节点
	err = stub.DelState(forfeitureAddr)
	if err != nil {
		return err
	}
	return nil
}

//处理没收Mediator保证金
func (d *DepositChaincode) handleMediatorForfeitureDeposit(foundationAddr string, forfeiture *modules.Forfeiture, balance *modules.DepositBalance, stub shim.ChaincodeStubInterface) pb.Response {
	var err error
	//计算余额
	result := balance.TotalAmount - forfeiture.ApplyTokens.Amount
	//判断是否没收全部，即在列表中移除该节点
	if result == 0 {
		//没收不考虑是否在规定周期内,其实它肯定是在列表中并已在周期内
		//没收全部，即删除,已经是计算好奖励了
		err = d.forfeitureAllDeposit("MediatorList", stub, foundationAddr, forfeiture.ForfeitureAddress, forfeiture.ApplyTokens)
		if err != nil {
			return shim.Success([]byte(err.Error()))
		}
		return shim.Success([]byte("成功退出"))
	} else {
		//TODO 对于mediator，要么全没收，要么退出一部分，且退出该部分金额后还在列表中
		return d.forfeitureSomeDeposit("Mediator", stub, foundationAddr, forfeiture, balance)
	}
}

func (d *DepositChaincode) forfertureAndMoveList(role string, stub shim.ChaincodeStubInterface, foundationAddr string, forfeiture *modules.Forfeiture, balance *modules.DepositBalance) pb.Response {
	//调用从合约把token转到请求地址
	err := stub.PayOutToken(foundationAddr, forfeiture.ApplyTokens, 0)
	if err != nil {
		return shim.Error(err.Error())
	}
	//handleMember(role, forfeiture.ForfeitureAddress, stub)
	err = moveCandidate(role, forfeiture.ForfeitureAddress, stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	//计算一部分的利息
	//获取币龄
	awards := award.GetAwardsWithCoins(balance.TotalAmount, balance.LastModifyTime.Unix())
	fmt.Println("awards ", awards)
	balance.LastModifyTime = time.Now().UTC()
	//加上利息奖励
	balance.TotalAmount += awards
	//减去提取部分
	balance.TotalAmount -= forfeiture.ApplyTokens.Amount

	balance.ForfeitureValues = append(balance.ForfeitureValues, forfeiture)

	//序列化
	return d.marshalForBalance(stub, forfeiture.ForfeitureAddress, balance)
}

//不需要移除候选列表，但是要没收一部分保证金
func (d *DepositChaincode) forfeitureSomeDeposit(role string, stub shim.ChaincodeStubInterface, foundationAddr string, forfeiture *modules.Forfeiture, balance *modules.DepositBalance) pb.Response {
	//调用从合约把token转到请求地址
	err := stub.PayOutToken(foundationAddr, forfeiture.ApplyTokens, 0)
	if err != nil {
		return shim.Error(err.Error())
	}
	//计算当前币龄奖励
	awards := award.GetAwardsWithCoins(balance.TotalAmount, balance.LastModifyTime.Unix())
	fmt.Println("awards ", awards)
	balance.LastModifyTime = time.Now().UTC()
	//加上利息奖励
	balance.TotalAmount += awards
	//减去提取部分
	balance.TotalAmount -= forfeiture.ApplyTokens.Amount

	balance.ForfeitureValues = append(balance.ForfeitureValues, forfeiture)

	//序列化
	return d.marshalForBalance(stub, forfeiture.ForfeitureAddress, balance)
}

func (d *DepositChaincode) handleJuryForfeitureDeposit(foundationAddr string, forfeiture *modules.Forfeiture, balance *modules.DepositBalance, stub shim.ChaincodeStubInterface) pb.Response {
	var err error
	//计算余额
	result := balance.TotalAmount - forfeiture.ApplyTokens.Amount
	//判断是否没收全部，即在列表中移除该节点
	if result == 0 {
		//没收不考虑是否在规定周期内,其实它肯定是在列表中并已在周期内
		//没收全部，即删除
		err = d.forfeitureAllDeposit("JuryList", stub, foundationAddr, forfeiture.ForfeitureAddress, forfeiture.ApplyTokens)
		if err != nil {
			return shim.Success([]byte(err.Error()))
		}
		return shim.Success([]byte("成功退出"))

	} else if result < depositAmountsForJury {
		//TODO 对于jury，需要移除列表
		return d.forfertureAndMoveList("JuryList", stub, foundationAddr, forfeiture, balance)
	} else {
		//TODO 退出一部分，且退出该部分金额后还在列表中
		return d.forfeitureSomeDeposit("Jury", stub, foundationAddr, forfeiture, balance)
	}
}

func (d *DepositChaincode) handleDeveloperForfeitureDeposit(foundationAddr string, forfeiture *modules.Forfeiture, balance *modules.DepositBalance, stub shim.ChaincodeStubInterface) pb.Response {
	var err error
	//计算余额
	result := balance.TotalAmount - forfeiture.ApplyTokens.Amount
	//判断是否没收全部，即在列表中移除该节点
	if result == 0 {
		//没收不考虑是否在规定周期内,其实它肯定是在列表中并已在周期内
		//没收全部，即删除
		err = d.forfeitureAllDeposit("DeveloperList", stub, foundationAddr, forfeiture.ForfeitureAddress, forfeiture.ApplyTokens)
		if err != nil {
			return shim.Success([]byte(err.Error()))
		}
		return shim.Success([]byte("成功退出"))
	} else if result < depositAmountsForDeveloper {
		return d.forfertureAndMoveList("DeveloperList", stub, foundationAddr, forfeiture, balance)
	} else {
		//TODO 退出一部分，且退出该部分金额后还在列表中
		return d.forfeitureSomeDeposit("Developer", stub, foundationAddr, forfeiture, balance)
	}
}
