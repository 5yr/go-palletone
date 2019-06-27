/*
 *
 *    This file is part of go-palletone.
 *    go-palletone is free software: you can redistribute it and/or modify
 *    it under the terms of the GNU General Public License as published by
 *    the Free Software Foundation, either version 3 of the License, or
 *    (at your option) any later version.
 *    go-palletone is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    GNU General Public License for more details.
 *    You should have received a copy of the GNU General Public License
 *    along with go-palletone.  If not, see <http://www.gnu.org/licenses/>.
 * /
 *
 *  * @author PalletOne core developer <dev@pallet.one>
 *  * @date 2018-2019
 *
 */

package modules

import "errors"

//质押列表
type PledgeList struct {
	TotalAmount uint64           `json:"total_amount"`
	Date        string           `json:"date"` //质押列表所在的日期yyyyMMdd
	Members     []*AddressAmount `json:"members"`
}

//账户质押情况
type AddressAmount struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
}

func NewAddressAmount(addr string, amt uint64) *AddressAmount {
	return &AddressAmount{Address: addr, Amount: amt}
}
func (pl *PledgeList) Add(addr string, amount uint64) {
	pl.TotalAmount += amount
	for _, p := range pl.Members {
		if p.Address == addr {
			p.Amount += amount
			return
		}
	}
	pl.Members = append(pl.Members, &AddressAmount{Address: addr, Amount: amount})
}

//从质押列表中提币，Amount 0表示全部提取
func (pl *PledgeList) Reduce(addr string, amount uint64) (uint64, error) {
	for i, p := range pl.Members {
		if p.Address == addr {
			if amount == 0 {
				amount = p.Amount //如果是0表示全部提取
			}
			if p.Amount < amount {
				return 0, errors.New("Not enough amount")
			}
			pl.TotalAmount -= amount
			if p.Amount == amount {
				pl.Members = append(pl.Members[:i], pl.Members[i+1:]...)
				return amount, nil
			}
			p.Amount -= amount
			return amount, nil
		}
	}
	return 0, errors.New("Address not found")
}
