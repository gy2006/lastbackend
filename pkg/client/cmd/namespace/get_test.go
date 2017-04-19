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

package namespace_test

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/client/cmd/namespace"
	"github.com/lastbackend/lastbackend/pkg/client/context"
	n "github.com/lastbackend/lastbackend/pkg/daemon/namespace/views/v1"
	h "github.com/lastbackend/lastbackend/pkg/util/http"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet(t *testing.T) {

	const (
		tName string = "test name"
		tDesc        = "test description"
	)

	var (
		err error
		ctx = context.Mock()
	)

	//------------------------------------------------------------------------------------------
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nspaceJSON, err := json.Marshal(n.Namespace{
			Meta: n.NamespaceMeta{
				Name:        tName,
				Description: tDesc,
			},
		})
		assert.NoError(t, err)

		w.WriteHeader(200)
		_, err = w.Write(nspaceJSON)
		assert.NoError(t, err)
	}))
	defer server.Close()
	//------------------------------------------------------------------------------------------

	ctx.SetHttpClient(h.New(server.URL[7:]))

	nspace, err := namespace.Get(tName)
	assert.NoError(t, err)
	assert.Equal(t, tName, nspace.Meta.Name)
	assert.Equal(t, tDesc, nspace.Meta.Description)
}
