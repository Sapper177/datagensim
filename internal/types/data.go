package types

import (
    "gosim/internal/engine";
)

type DataItem struct {
    Name      string

}




func (d *DataItem) Update() {
    switch d.Engine {
    case Static:
        // do nothing
    case Ramp:
        d.Value += d.Step
        if d.Value > d.Max {
            d.Value = d.Min
        }
    case Sinusoidal:
        d.counter += 1
        d.Value = d.Min + (d.Max-d.Min)/2 + (d.Max-d.Min)/2*math.Sin(2*math.Pi*d.Frequency*d.counter+d.Phase)
    }
}