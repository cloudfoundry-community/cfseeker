package api

import (
	"net/http"

	"github.com/cloudfoundry-community/cfseeker/seeker"
)

func invalidateBOSHCacheHandler(w http.ResponseWriter, r *http.Request, s *seeker.Seeker) {
	s.InvalidateAll()
	NewResponse(w).Message("BOSH VM info cache successfully cleared").Write()
}
