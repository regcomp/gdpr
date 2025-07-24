package caching

// func (ns *NonceStore) cleanupDaemon() {
// 	ticker := time.NewTicker(nonceCleanupInterval * time.Minute)
// 	defer ticker.Stop()
//
// 	for {
// 		select {
// 		case <-ns.ctx.Done():
// 			return
// 		case <-ticker.C:
// 			ns.cleanup()
// 		}
// 	}
// }
//
// func (ns *NonceStore) cleanup() {
// 	ns.mu.RLock()
// 	now := time.Now()
// 	toDelete := make([]string, 0)
// 	for nonce, expiry := range ns.nonces {
// 		if now.After(expiry) {
// 			toDelete = append(toDelete, nonce)
// 		}
// 	}
// 	ns.mu.RUnlock()
//
// 	if len(toDelete) > 0 {
// 		ns.mu.Lock()
// 		for _, nonce := range toDelete {
// 			delete(ns.nonces, nonce)
// 		}
// 		ns.mu.Unlock()
// 	}
// }
