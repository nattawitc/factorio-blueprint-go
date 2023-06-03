package factoriobp

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"strings"
)

type object struct {
	Blueprint     *Blueprint     `json:"blueprint,omitempty"`
	BlueprintBook *BlueprintBook `json:"blueprint_book,omitempty"`
}

func decode(s string) (object, error) {
	r, err := zlib.NewReader(base64.NewDecoder(base64.StdEncoding, strings.NewReader(s)))
	if err != nil {
		return object{}, err
	}

	var obj object
	if err := json.NewDecoder(r).Decode(&obj); err != nil {
		return object{}, err
	}
	return obj, nil
}

func (o object) encode() (string, error) {
	var out bytes.Buffer
	out.Write([]byte("0"))

	w, err := zlib.NewWriterLevel(base64.NewEncoder(base64.RawStdEncoding, &out), zlib.BestCompression)
	if err != nil {
		return "", err
	}
	defer w.Close()
	defer w.Flush()

	if err := json.NewEncoder(w).Encode(o); err != nil {
		return "", err
	}

	w.Close()

	return out.String(), nil
}

type BlueprintBook struct {
	Item        string      `json:"item"`
	Label       string      `json:"label"`
	LabelColor  *Color      `json:"label_color,omitempty"`
	Blueprints  []Blueprint `json:"blueprints"`
	ActiveIndex int         `json:"active_index"`
	Version     uint64      `json:"version"`
}

func DecodeBlueprintBook(s string) (BlueprintBook, error) {
	obj, err := decode(s)
	if err != nil || obj.BlueprintBook == nil {
		return BlueprintBook{}, err
	}
	return *obj.BlueprintBook, nil
}

func (bb BlueprintBook) Encode() (string, error) {
	obj := object{
		BlueprintBook: &bb,
	}

	return obj.encode()
}

type Blueprint struct {
	Item       string     `json:"item"`
	Label      string     `json:"label"`
	LabelColor *Color     `json:"label_color,omitempty"`
	Entities   []Entity   `json:"entities"`
	Tiles      []Tile     `json:"tiles"`
	Icons      []Icon     `json:"icons"`
	Schedules  []Schedule `json:"schedules"`
	Version    int64      `json:"version"`
}

func DecodeBlueprint(s string) (Blueprint, error) {
	obj, err := decode(s)
	if err != nil || obj.Blueprint == nil {
		return Blueprint{}, err
	}
	return *obj.Blueprint, nil
}

func (bp Blueprint) Encode() (string, error) {
	obj := object{
		Blueprint: &bp,
	}

	return obj.encode()
}

type Icon struct {
	Index  int      `json:"index"`
	Signal SignalID `json:"signal"`
}

type SignalID struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Entity struct {
	EntityNumber       int                    `json:"entity_number"`
	Name               string                 `json:"name"`
	Position           Position               `json:"position"`
	Direction          int                    `json:"direction,omitempty"`
	Orientation        float64                `json:"orientation,omitempty"`
	Connectgions       *Connection            `json:"connectgions,omitempty"`
	Neighbours         []int                  `json:"neighbours,omitempty"`
	ControlBehavior    map[string]any         `json:"control_behavior,omitempty"`
	Items              map[string]uint        `json:"items,omitempty"`
	Recipe             string                 `json:"recipe,omitempty"`
	Bar                int                    `json:"bar,omitempty"`
	Inventory          *Inventory             `json:"inventory,omitempty"`
	InfinitySettings   *InfinitySetting       `json:"infinity_settings,omitempty"`
	Type               string                 `json:"type,omitempty"`
	InputPriority      string                 `json:"input_priority,omitempty"`
	OutputPriority     string                 `json:"output_priority,omitempty"`
	Filter             string                 `json:"filter,omitempty"`
	Filters            []ItemFilter           `json:"filters,omitempty"`
	FilterMode         string                 `json:"filter_mode,omitempty"`
	OverrideStackSize  uint                   `json:"override_stack_size,omitempty"`
	DropPosition       *Position              `json:"drop_position,omitempty"`
	PickupPosition     *Position              `json:"pickup_position,omitempty"`
	RequestFilters     *LogisticFilter        `json:"request_filters,omitempty"`
	RequestFromBuffers bool                   `json:"request_from_buffers,omitempty"`
	Parameters         *SpeakerParameter      `json:"parameters,omitempty"`
	AlertMessage       *SpeakerAlertParameter `json:"alert_message,omitempty"`
	AutoLaunch         bool                   `json:"auto_launch,omitempty"`
	Variation          uint                   `json:"variation,omitempty"`
	Color              *Color                 `json:"color,omitempty"`
	Station            string                 `json:"station,omitempty"`
}

type Inventory struct {
	Filters []ItemFilter `json:"filters"`
	Bar     *int         `json:"bar,omitempty"`
}

type Schedule struct {
	Schedule    []ScheduleRecord `json:"schedule"`
	Locomotives []int            `json:"locomotives"`
}

type ScheduleRecord struct {
	Station        string          `json:"station"`
	WaitConditions []WaitCondition `json:"wait_conditions,omitempty"`
}

type WaitCondition struct {
	Type        string         `json:"type"`
	CompareType string         `json:"compare_type"`
	Ticks       uint           `json:"ticks,omitempty"`
	Condition   map[string]any `json:"condition,omitempty"`
}

type Tile struct {
	Name     string   `json:"name"`
	Position Position `json:"position"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Connection struct {
	One ConnectionPoint `json:"1"`
	Two ConnectionPoint `json:"2"`
}

type ConnectionPoint struct {
	Red   []ConnectionData `json:"red"`
	Green []ConnectionData `json:"green"`
}

type ConnectionData struct {
	EntityID int `json:"entity_id"`
	Circuit  int `json:"circuit"`
}

type ItemFilter struct {
	Name  string `json:"name"`
	Index int    `json:"index"`
}

type InfinitySetting struct {
	RemoveUnfilteredItems bool             `json:"remove_unfiltered_items"`
	Filters               []InfinityFilter `json:"filters"`
}

type InfinityFilter struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Mode  string `json:"mode"`
	Index int    `json:"index"`
}

type LogisticFilter struct {
	Name  string `json:"name"`
	Index int    `json:"index"`
	Count int    `json:"count"`
}

type SpeakerParameter struct {
	PlaybackVolume   float64 `json:"playback_volume"`
	PlaybackGlobally bool    `json:"playback_globally"`
	AllowPolyphony   bool    `json:"allow_polyphony"`
}

type SpeakerAlertParameter struct {
	ShowAlert    bool     `json:"show_alert"`
	ShowOnMap    bool     `json:"show_on_map"`
	IconSignalID SignalID `json:"icon_signal_id"`
	AlertMessage string   `json:"alert_message"`
}

type Color struct {
	R float64 `json:"r"`
	G float64 `json:"g"`
	B float64 `json:"b"`
	A float64 `json:"a"`
}
