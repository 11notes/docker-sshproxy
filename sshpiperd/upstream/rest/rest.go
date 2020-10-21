package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type restPipe struct {
	Direction string `json:"Direction"`
	User string `json:"User"`
}

func (p *plugin) findUpstream(conn ssh.ConnMetadata, challengeContext ssh.AdditionalChallengeContext) (net.Conn, *ssh.AuthPipe, error) {

	user := conn.User()
	remoteaddr := conn.RemoteAddr()
	localaddr := conn.LocalAddr()
	d, err := getUpstream(user)

	if err != nil {
		return nil, nil, err
	}

	addr := d.Upstream.Server.Address
	upuser := d.Upstream.Username

	if upuser == "" {
		upuser = d.Username
	}

	logger.Printf("mapping downstream [%v@%v from %v] to upstream [%v@%v]", user, localaddr, remoteaddr, upuser, addr)

	c, err := upstreamprovider.DialForSSH(addr)

	if err != nil {
		return nil, nil, err
	}

	hostKeyCallback := ssh.InsecureIgnoreHostKey()

	if !d.Upstream.Server.IgnoreHostKey {

		key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(d.Upstream.Server.HostKey.Key.Data))
		if err != nil {
			return nil, nil, err
		}

		hostKeyCallback = ssh.FixedHostKey(key)
	}

	pipe := ssh.AuthPipe{
		User: upuser,

		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (ssh.AuthPipeType, ssh.AuthMethod, error) {

			expectKey := key.Marshal()
			for _, k := range d.AuthorizedKeys {
				publicKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(k.Key.Data))

				if err != nil {
					logger.Printf("parse [keyid = %v] error :%v. skip to next key", k.Key.ID, err)
					continue
				}

				if bytes.Equal(publicKey.Marshal(), expectKey) {

					kinterf, err := ssh.ParseRawPrivateKey([]byte(d.Upstream.PrivateKey.Key.Data))
					if err != nil {
						break
					}

					signer, err := ssh.NewSignerFromKey(kinterf)
					if err != nil || signer == nil {
						break
					}

					return ssh.AuthPipeTypeMap, ssh.PublicKeys(signer), nil
				}
			}

			return ssh.AuthPipeTypeNone, nil, nil
		},

		UpstreamHostKeyCallback: hostKeyCallback,
	}

	return c, &pipe, nil
}

func (p *rest) getDownstream(User string) (*pipe, error) {
	jsonStr, err := json.Marshal(restPipe{Direction:"Downstream", User:User})
    req, err := http.NewRequest("POST", p.Config.URL, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad http state code %v", resp.StatusCode)
	}

    body, _ := ioutil.ReadAll(resp.Body)
    if err != nil {
		return nil, err
	}

	pipe := pipe{}

	if err := json.Unmarshal(body, &pipe); err != nil {
		return nil, err
	}

	return &pipe, nil
}

func (p *rest) getUpstream(User string) (*pipe, error) {
	jsonStr, err := json.Marshal(restPipe{Direction:"Upstream", User:User})
    req, err := http.NewRequest("POST", p.Config.URL, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad http state code %v", resp.StatusCode)
	}

    body, _ := ioutil.ReadAll(resp.Body)
    if err != nil {
		return nil, err
	}

	pipe := pipe{}

	if err := json.Unmarshal(body, &pipe); err != nil {
		return nil, err
	}

	return &pipe, nil
}