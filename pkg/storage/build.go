//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package storage

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const buildStorage string = "builds"

// Service Build type for interface in interfaces folder
type BuildStorage struct {
	IBuild
	util   IUtil
	Client func() (store.IStore, store.DestroyFunc, error)
}

// Get build model by id
func (s *BuildStorage) GetByID(ctx context.Context, id string) (*types.Build, error) {
	return nil, nil
}

// Get builds by image
func (s *BuildStorage) ListByImage(ctx context.Context, id string) (*types.BuildList, error) {
	return nil, nil
}

// Insert new build into storage
func (s *BuildStorage) Insert(ctx context.Context, build *types.Build) (*types.Build, error) {
	return nil, nil
}

func NewBuildStorage(config store.Config, util IUtil) *BuildStorage {
	s := new(BuildStorage)
	s.util = util
	s.Client = func() (store.IStore, store.DestroyFunc, error) {
		return New(config)
	}
	return s
}
