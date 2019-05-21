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
 *
 */

package storage

import (
	"github.com/palletone/go-palletone/common/log"
	"github.com/palletone/go-palletone/common/ptndb"
	"github.com/palletone/go-palletone/dag/constants"
	"github.com/palletone/go-palletone/dag/modules"
)

var (
	GlobalPropDBKey    = append(constants.GLOBALPROPERTY_PREFIX, []byte("GlobalProperty")...)
	DynGlobalPropDBKey = append(constants.DYNAMIC_GLOBALPROPERTY_PREFIX, []byte("DynamicGlobalProperty")...)
)

func StoreGlobalProp(db ptndb.Database, gp *modules.GlobalProperty) error {
	err := StoreToJsonBytes(db, GlobalPropDBKey, gp)
	if err != nil {
		log.Errorf("Store global properties error: %v", err.Error())
	}

	return err
}

func RetrieveGlobalProp(db ptndb.Database) (*modules.GlobalProperty, error) {
	gp := modules.NewGlobalProp()

	err := RetrieveFromJsonBytes(db, GlobalPropDBKey, gp)
	if err != nil {
		log.Errorf("Retrieve global properties error: %v", err.Error())
	}

	return gp, err
}

func StoreDynGlobalProp(db ptndb.Database, dgp *modules.DynamicGlobalProperty) error {
	err := StoreBytes(db, DynGlobalPropDBKey, *dgp)
	if err != nil {
		log.Errorf("Store dynamic global properties error: %v", err.Error())
	}

	return err
}

func RetrieveDynGlobalProp(db ptndb.Database) (*modules.DynamicGlobalProperty, error) {
	dgp := modules.NewDynGlobalProp()

	err := Retrieve(db, DynGlobalPropDBKey, dgp)
	if err != nil {
		log.Errorf("Retrieve dynamic global properties error: %v", err.Error())
	}

	return dgp, err
}
