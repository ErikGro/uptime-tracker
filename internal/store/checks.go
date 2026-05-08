package store

import "time"

func (s *Store) AppendCheck(c *Check) error {
	return s.db.Create(c).Error
}

func (s *Store) ListChecksFor(urlID uint, limit int) ([]Check, error) {
	var checks []Check
	err := s.db.Where("url_id = ?", urlID).
		Order("checked_at DESC").
		Limit(limit).
		Find(&checks).Error
	return checks, err
}

func (s *Store) PruneChecksOlderThan(cutoff time.Time) (int64, error) {
	res := s.db.Where("checked_at < ?", cutoff).Delete(&Check{})
	return res.RowsAffected, res.Error
}
