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
 * @author PalletOne core developer Albert·Gou <dev@pallet.one>
 * @date 2018
 */

package mediatorplugin

import (
	"fmt"

	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/event"
	"github.com/palletone/go-palletone/common/log"
	"go.dedis.ch/kyber/v3/share/dkg/pedersen"
	"go.dedis.ch/kyber/v3/share/vss/pedersen"
)

func (mp *MediatorPlugin) startVSSProtocol() {
	log.Debugf("Start completing the VSS protocol.")

	// 处理其他 mediator 的 deals
	go mp.launchProcessDealLoops()

	// todo albert 处理所有的response

	// todo albert 换届后第一个生产槽的前一个生产间隔开始vss协议

	// 将deal广播给其他节点
	go mp.broadcastVSSDeals()

	// 开启计时器，删除vss相关缓存
}

func (mp *MediatorPlugin) launchProcessDealLoops() {
	lams := mp.GetLocalActiveMediators()

	for _, localMed := range lams {
		go mp.processDealLoop(localMed)
	}
}

func (mp *MediatorPlugin) processDealLoop(localMed common.Address) {

}

func (mp *MediatorPlugin) processVSSDeal(dealEvent *VSSDealEvent) {
	dag := mp.dag
	localMed := dag.GetActiveMediatorAddr(int(dealEvent.DstIndex))

	dkgr, err := mp.getLocalActiveDKG(localMed)
	if err != nil {
		log.Debugf(err.Error())
		return
	}

	deal := dealEvent.Deal

	resp, err := dkgr.ProcessDeal(deal)
	if err != nil {
		log.Debugf("dkg: cannot process own deal: " + err.Error())
		return
	}

	vrfrMed := dag.GetActiveMediatorAddr(int(deal.Index))
	log.Debugf("the mediator(%v) received the vss deal from the mediator(%v)",
		localMed.Str(), vrfrMed.Str())

	// todo albert 待重构
	go mp.processResponseLoop(localMed, vrfrMed)

	if resp.Response.Status != vss.StatusApproval {
		err = fmt.Errorf("DKG: own deal gave a complaint: %v", localMed.String())
		log.Debugf(err.Error())
		return
	}

	respEvent := VSSResponseEvent{
		Resp: resp,
	}
	go mp.vssResponseFeed.Send(respEvent)
	log.Debugf("the mediator(%v) broadcast the vss response to the mediator(%v)",
		localMed.Str(), vrfrMed.Str())

	return
}

func (mp *MediatorPlugin) broadcastVSSDeals() {
	for localMed, dkg := range mp.activeDKGs {
		deals, err := dkg.Deals()
		if err != nil {
			log.Debugf(err.Error())
			continue
		}

		// todo albert 待重构
		go mp.processResponseLoop(localMed, localMed)
		log.Debugf("the mediator(%v) broadcast vss deals", localMed.Str())

		for index, deal := range deals {
			event := VSSDealEvent{
				DstIndex: uint32(index),
				Deal:     deal,
			}

			go mp.vssDealFeed.Send(event)
		}
	}
}

func (mp *MediatorPlugin) SubscribeVSSDealEvent(ch chan<- VSSDealEvent) event.Subscription {
	return mp.vssDealScope.Track(mp.vssDealFeed.Subscribe(ch))
}

func (mp *MediatorPlugin) AddToDealBuf(dealEvent *VSSDealEvent) {
	if !mp.groupSigningEnabled {
		return
	}

	dag := mp.dag
	localMed := dag.GetActiveMediatorAddr(int(dealEvent.DstIndex))

	deal := dealEvent.Deal
	vrfrMed := dag.GetActiveMediatorAddr(int(deal.Index))
	log.Debugf("the mediator(%v) received the vss deal from the mediator(%v)",
		localMed.Str(), vrfrMed.Str())

	if _, ok := mp.dealBuf[localMed]; !ok {
		aSize := dag.ActiveMediatorsCount()
		mp.dealBuf[localMed] = make(chan *dkg.Deal, aSize-1)
	}

	mp.dealBuf[localMed] <- deal
}

func (mp *MediatorPlugin) AddToResponseBuf(respEvent *VSSResponseEvent) {
	if !mp.groupSigningEnabled {
		return
	}

	resp := respEvent.Resp
	lams := mp.GetLocalActiveMediators()
	for _, localMed := range lams {
		dag := mp.dag

		//ignore the message from myself
		srcIndex := resp.Response.Index
		srcMed := dag.GetActiveMediatorAddr(int(srcIndex))
		if srcMed == localMed {
			continue
		}

		vrfrMed := dag.GetActiveMediatorAddr(int(resp.Index))
		log.Debugf("the mediator(%v) received the vss response from the mediator(%v) to the mediator(%v)",
			localMed.Str(), srcMed.Str(), vrfrMed.Str())

		if _, ok := mp.respBuf[localMed][vrfrMed]; !ok {
			log.Debugf("the mediator(%v)'s respBuf corresponding the mediator(%v) is not initialized",
				localMed.Str(), vrfrMed.Str())
		}
		mp.respBuf[localMed][vrfrMed] <- resp
	}
}

func (mp *MediatorPlugin) SubscribeVSSResponseEvent(ch chan<- VSSResponseEvent) event.Subscription {
	return mp.vssResponseScope.Track(mp.vssResponseFeed.Subscribe(ch))
}

func (mp *MediatorPlugin) processResponseLoop(localMed, vrfrMed common.Address) {
	dkgr, err := mp.getLocalActiveDKG(localMed)
	if err != nil {
		log.Debugf(err.Error())
		return
	}

	respAmount := mp.dag.ActiveMediatorsCount() - 1
	respCount := 0
	// localMed 对 vrfrMed 的 response 在 ProcessDeal 生成 response 时 自动处理了
	if vrfrMed != localMed {
		respCount++
	}

	processResp := func(resp *dkg.Response) bool {
		jstf, err := dkgr.ProcessResponse(resp)
		if err != nil {
			log.Debugf(err.Error())
			return false
		}

		if jstf != nil {
			log.Debugf("DKG: wrong Process Response: %v", localMed.String())
			return false
		}

		return true
	}

	isFinishedAndCertified := func() (finished, certified bool) {
		respCount++

		if respCount == respAmount-1 {
			finished = true

			if dkgr.Certified() {
				log.Debugf("the mediator(%v)'s DKG verification passed", localMed.Str())

				certified = true
			}
		}

		return
	}

	if _, ok := mp.respBuf[localMed][vrfrMed]; !ok {
		log.Debugf("the mediator(%v)'s respBuf corresponding the mediator(%v) is not initialized",
			localMed.Str(), vrfrMed.Str())
	}

	log.Debugf("the mediator(%v) run the loop to process response regarding the mediator(%v)",
		localMed.Str(), vrfrMed.Str())
	respCh := mp.respBuf[localMed][vrfrMed]

	for {
		select {
		case <-mp.quit:
			return
			// todo Albert 超过期限后也要删除
		case resp := <-respCh:
			processResp(resp)
			finished, certified := isFinishedAndCertified()
			if finished {
				delete(mp.respBuf[localMed], vrfrMed)

				if certified {
					go mp.signUnitsTBLS(localMed)
					go mp.recoverUnitsTBLS(localMed)

					//delete(mp.respBuf, localMed)
				}

				return
			}
		}
	}
}
