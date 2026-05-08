package store

import (
	"strconv"
	"time"

	"gorm.io/gorm/clause"
)

const (
	KeyPollInterval     = "poll_interval_seconds"
	KeyRequestTimeout   = "request_timeout_seconds"
	KeyFailureThreshold = "failure_threshold"
	KeyWebhookURL       = "webhook_url"
	KeyWebhookEnabled   = "webhook_enabled"
	KeyRetentionDays    = "retention_days"
)

var settingDefaults = map[string]string{
	KeyPollInterval:     "300",
	KeyRequestTimeout:   "10",
	KeyFailureThreshold: "3",
	KeyWebhookURL:       "",
	KeyWebhookEnabled:   "false",
	KeyRetentionDays:    "30",
}

func (s *Store) seedSettings() error {
	rows := make([]Setting, 0, len(settingDefaults))
	for k, v := range settingDefaults {
		rows = append(rows, Setting{Key: k, Value: v})
	}
	return s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&rows).Error
}

func (s *Store) GetSetting(key string) (string, error) {
	var setting Setting
	if err := s.db.First(&setting, "key = ?", key).Error; err != nil {
		if isNotFound(err) {
			if def, ok := settingDefaults[key]; ok {
				return def, nil
			}
		}
		return "", err
	}
	return setting.Value, nil
}

func (s *Store) SetSetting(key, value string) error {
	return s.db.Save(&Setting{Key: key, Value: value}).Error
}

func (s *Store) AllSettings() (map[string]string, error) {
	var rows []Setting
	if err := s.db.Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]string, len(rows))
	for _, row := range rows {
		out[row.Key] = row.Value
	}
	return out, nil
}

func (s *Store) PollInterval() time.Duration {
	return time.Duration(s.intSetting(KeyPollInterval)) * time.Second
}

func (s *Store) RequestTimeout() time.Duration {
	return time.Duration(s.intSetting(KeyRequestTimeout)) * time.Second
}

func (s *Store) FailureThreshold() int {
	return s.intSetting(KeyFailureThreshold)
}

func (s *Store) WebhookURL() string {
	v, _ := s.GetSetting(KeyWebhookURL)
	return v
}

func (s *Store) WebhookEnabled() bool {
	v, _ := s.GetSetting(KeyWebhookEnabled)
	b, _ := strconv.ParseBool(v)
	return b
}

func (s *Store) RetentionDays() int {
	return s.intSetting(KeyRetentionDays)
}

func (s *Store) intSetting(key string) int {
	v, err := s.GetSetting(key)
	if err != nil {
		return 0
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return n
}
