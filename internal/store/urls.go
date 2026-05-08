package store

import "database/sql"

func (s *Store) ListURLs() ([]URL, error) {
	var urls []URL
	err := s.db.Order("created_at DESC").Find(&urls).Error
	return urls, err
}

func (s *Store) GetURL(id uint) (*URL, error) {
	var u URL
	if err := s.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) CreateURL(label, target string) (*URL, error) {
	u := &URL{Label: label, URL: target, CurrentStatus: StatusUnknown}
	if err := s.db.Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Store) UpdateURL(id uint, label, target string) (*URL, error) {
	var u URL
	if err := s.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	u.Label = label
	u.URL = target
	if err := s.db.Save(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) DeleteURL(id uint) error {
	return s.db.Delete(&URL{}, id).Error
}

func (s *Store) UpdateURLStatus(id uint, status Status, consecutiveFailures int, lastCheckedAt sql.NullTime) error {
	return s.db.Model(&URL{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"current_status":       status,
			"consecutive_failures": consecutiveFailures,
			"last_checked_at":      lastCheckedAt,
		}).Error
}
