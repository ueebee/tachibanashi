package master

import (
	"sort"
	"strings"
	"sync"

	"github.com/ueebee/tachibanashi/model"
)

type MasterType string

const (
	MasterSystemStatus          MasterType = "CLMSystemStatus"
	MasterDateZyouhou           MasterType = "CLMDateZyouhou"
	MasterYobine                MasterType = "CLMYobine"
	MasterUnyouStatus           MasterType = "CLMUnyouStatus"
	MasterUnyouStatusKabu       MasterType = "CLMUnyouStatusKabu"
	MasterUnyouStatusHasei      MasterType = "CLMUnyouStatusHasei"
	MasterIssueMstKabu          MasterType = "CLMIssueMstKabu"
	MasterIssueSizyouMstKabu    MasterType = "CLMIssueSizyouMstKabu"
	MasterIssueSizyouKiseiKabu  MasterType = "CLMIssueSizyouKiseiKabu"
	MasterIssueMstSak           MasterType = "CLMIssueMstSak"
	MasterIssueMstOp            MasterType = "CLMIssueMstOp"
	MasterIssueSizyouKiseiHasei MasterType = "CLMIssueSizyouKiseiHasei"
	MasterDaiyouKakeme          MasterType = "CLMDaiyouKakeme"
	MasterHosyoukinMst          MasterType = "CLMHosyoukinMst"
	MasterOrderErrReason        MasterType = "CLMOrderErrReason"
	MasterEventDownloadComplete MasterType = "CLMEventDownloadComplete"
	MasterIssueMstOther         MasterType = "CLMIssueMstOther"
	MasterIssueMstIndex         MasterType = "CLMIssueMstIndex"
	MasterIssueMstFx            MasterType = "CLMIssueMstFx"
)

type UpdateMeta struct {
	Serial    int64
	UpdatedAt string
	Deleted   bool
}

type Record struct {
	Key    string
	Fields model.Attributes
	Meta   UpdateMeta
}

type MasterStore interface {
	Get(typ MasterType, key string) (Record, bool)
	Upsert(typ MasterType, key string, fields model.Attributes, meta UpdateMeta) bool
	Delete(typ MasterType, key string) bool
	All(typ MasterType) []Record
}

type IndexSpec struct {
	Name   string
	Fields []string
}

const indexDelimiter = "\x1f"

func JoinIndex(parts ...string) string {
	return strings.Join(parts, indexDelimiter)
}

type MemoryStore struct {
	mu         sync.RWMutex
	records    map[MasterType]map[string]Record
	indexSpecs map[MasterType][]IndexSpec
	indexes    map[MasterType]map[string]map[string]map[string]struct{}
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		records:    make(map[MasterType]map[string]Record),
		indexSpecs: make(map[MasterType][]IndexSpec),
		indexes:    make(map[MasterType]map[string]map[string]map[string]struct{}),
	}
}

func (s *MemoryStore) RegisterIndex(typ MasterType, spec IndexSpec) {
	normalized, ok := normalizeIndexSpec(spec)
	if !ok {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existing := range s.indexSpecs[typ] {
		if existing.Name == normalized.Name {
			return
		}
	}
	s.indexSpecs[typ] = append(s.indexSpecs[typ], normalized)
	for _, record := range s.records[typ] {
		s.addIndexLocked(typ, normalized, record)
	}
}

func (s *MemoryStore) FindByIndex(typ MasterType, name, value string) []Record {
	name = strings.TrimSpace(name)
	if name == "" || strings.TrimSpace(value) == "" {
		return nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	typeIndexes := s.indexes[typ]
	if typeIndexes == nil {
		return nil
	}
	indexValues := typeIndexes[name]
	if indexValues == nil {
		return nil
	}
	keys := indexValues[value]
	if len(keys) == 0 {
		return nil
	}
	out := make([]Record, 0, len(keys))
	keyList := make([]string, 0, len(keys))
	for key := range keys {
		keyList = append(keyList, key)
	}
	sort.Strings(keyList)
	for _, key := range keyList {
		record, ok := s.records[typ][key]
		if !ok {
			continue
		}
		out = append(out, cloneRecord(record))
	}
	return out
}

func (s *MemoryStore) Get(typ MasterType, key string) (Record, bool) {
	key = strings.TrimSpace(key)
	if key == "" {
		return Record{}, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	record, ok := s.records[typ][key]
	if !ok {
		return Record{}, false
	}
	return cloneRecord(record), true
}

func (s *MemoryStore) All(typ MasterType) []Record {
	s.mu.RLock()
	defer s.mu.RUnlock()
	records := s.records[typ]
	if len(records) == 0 {
		return nil
	}
	keys := make([]string, 0, len(records))
	for key := range records {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]Record, 0, len(keys))
	for _, key := range keys {
		out = append(out, cloneRecord(records[key]))
	}
	return out
}

func (s *MemoryStore) Upsert(typ MasterType, key string, fields model.Attributes, meta UpdateMeta) bool {
	key = strings.TrimSpace(key)
	if key == "" {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if meta.Deleted {
		return s.deleteLocked(typ, key)
	}

	records := s.records[typ]
	if records == nil {
		records = make(map[string]Record)
		s.records[typ] = records
	}
	existing, ok := records[key]
	if ok && !shouldReplace(existing.Meta, meta) {
		return false
	}
	if ok {
		s.removeIndexesLocked(typ, existing)
	}
	record := Record{
		Key:    key,
		Fields: cloneAttributes(fields),
		Meta:   meta,
	}
	records[key] = record
	for _, spec := range s.indexSpecs[typ] {
		s.addIndexLocked(typ, spec, record)
	}
	return true
}

func (s *MemoryStore) Delete(typ MasterType, key string) bool {
	key = strings.TrimSpace(key)
	if key == "" {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.deleteLocked(typ, key)
}

func (s *MemoryStore) deleteLocked(typ MasterType, key string) bool {
	records := s.records[typ]
	if len(records) == 0 {
		return false
	}
	record, ok := records[key]
	if !ok {
		return false
	}
	delete(records, key)
	s.removeIndexesLocked(typ, record)
	return true
}

func (s *MemoryStore) addIndexLocked(typ MasterType, spec IndexSpec, record Record) {
	value, ok := buildIndexValue(spec, record.Fields)
	if !ok {
		return
	}
	typeIndexes := s.indexes[typ]
	if typeIndexes == nil {
		typeIndexes = make(map[string]map[string]map[string]struct{})
		s.indexes[typ] = typeIndexes
	}
	indexValues := typeIndexes[spec.Name]
	if indexValues == nil {
		indexValues = make(map[string]map[string]struct{})
		typeIndexes[spec.Name] = indexValues
	}
	keys := indexValues[value]
	if keys == nil {
		keys = make(map[string]struct{})
		indexValues[value] = keys
	}
	keys[record.Key] = struct{}{}
}

func (s *MemoryStore) removeIndexesLocked(typ MasterType, record Record) {
	typeIndexes := s.indexes[typ]
	if typeIndexes == nil {
		return
	}
	for _, spec := range s.indexSpecs[typ] {
		value, ok := buildIndexValue(spec, record.Fields)
		if !ok {
			continue
		}
		indexValues := typeIndexes[spec.Name]
		if indexValues == nil {
			continue
		}
		keys := indexValues[value]
		if keys == nil {
			continue
		}
		delete(keys, record.Key)
		if len(keys) == 0 {
			delete(indexValues, value)
		}
	}
}

func normalizeIndexSpec(spec IndexSpec) (IndexSpec, bool) {
	name := strings.TrimSpace(spec.Name)
	if name == "" {
		return IndexSpec{}, false
	}
	fields := make([]string, 0, len(spec.Fields))
	for _, field := range spec.Fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		fields = append(fields, field)
	}
	if len(fields) == 0 {
		return IndexSpec{}, false
	}
	return IndexSpec{Name: name, Fields: fields}, true
}

func buildIndexValue(spec IndexSpec, fields model.Attributes) (string, bool) {
	values := make([]string, 0, len(spec.Fields))
	for _, field := range spec.Fields {
		value := strings.TrimSpace(fields.Value(field))
		if value == "" {
			return "", false
		}
		values = append(values, value)
	}
	return JoinIndex(values...), true
}

func shouldReplace(existing, next UpdateMeta) bool {
	if existing.Serial > 0 || next.Serial > 0 {
		switch {
		case next.Serial > existing.Serial:
			return true
		case next.Serial < existing.Serial:
			return false
		}
	}
	if next.UpdatedAt != "" && existing.UpdatedAt != "" {
		return next.UpdatedAt >= existing.UpdatedAt
	}
	if next.UpdatedAt != "" && existing.UpdatedAt == "" {
		return true
	}
	if next.UpdatedAt == "" && existing.UpdatedAt != "" {
		return false
	}
	return true
}

func cloneRecord(record Record) Record {
	return Record{
		Key:    record.Key,
		Fields: cloneAttributes(record.Fields),
		Meta:   record.Meta,
	}
}

func cloneAttributes(attrs model.Attributes) model.Attributes {
	if attrs == nil {
		return nil
	}
	out := make(model.Attributes, len(attrs))
	for key, value := range attrs {
		out[key] = value
	}
	return out
}
