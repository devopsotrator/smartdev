package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/anyandrea/smartdev/lib/database/weatherdb"
	"github.com/anyandrea/smartdev/lib/web"
	"github.com/gorilla/mux"
)

func GetSensorValues(wdb weatherdb.WeatherDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error

		vars := mux.Vars(req)
		id := vars["id"]
		limit := vars["limit"]

		if len(id) > 0 {
			var sensorId, valueLimit int64

			if len(id) > 0 {
				sensorId, err = strconv.ParseInt(id, 10, 64)
				if err != nil {
					Error(rw, err)
					return
				}
			}
			if len(limit) > 0 {
				valueLimit, err = strconv.ParseInt(limit, 10, 64)
				if err != nil {
					Error(rw, err)
					return
				}
			} else {
				valueLimit = 100
			}

			values, err := wdb.GetSensorValues(int(sensorId), int(valueLimit))
			if err != nil {
				Error(rw, err)
				return
			}

			web.Render().JSON(rw, http.StatusOK, values)
			return
		}

		web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}

func AddSensorValue(wdb weatherdb.WeatherDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error

		vars := mux.Vars(req)
		id := vars["id"]

		if len(id) > 0 {
			var sensorId, value int64
			if len(id) > 0 {
				sensorId, err = strconv.ParseInt(id, 10, 64)
				if err != nil {
					Error(rw, err)
					return
				}
			}

			req.ParseForm()
			value, err = strconv.ParseInt(req.Form.Get("value"), 10, 64)
			if err != nil {
				Error(rw, err)
				return
			}

			if err := wdb.InsertSensorValue(int(sensorId), int(value), time.Now()); err != nil {
				Error(rw, err)
				return
			}

			web.Render().JSON(rw, http.StatusCreated, value)
			return
		}

		web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}

func DeleteSensorValues(wdb weatherdb.WeatherDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error

		vars := mux.Vars(req)
		id := vars["id"]

		if len(id) > 0 {
			var sensorId int64
			if len(id) > 0 {
				sensorId, err = strconv.ParseInt(id, 10, 64)
				if err != nil {
					Error(rw, err)
					return
				}
			}

			if err := wdb.DeleteSensorValues(int(sensorId)); err != nil {
				Error(rw, err)
				return
			}

			web.Render().JSON(rw, http.StatusNoContent, nil)
			return
		}

		web.Render().JSON(rw, http.StatusNotFound, nil)
	}
}

func Housekeeping(wdb weatherdb.WeatherDB) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		days, err := strconv.ParseInt(req.Form.Get("days"), 10, 64)
		if err != nil {
			Error(rw, err)
			return
		}
		rows, err := strconv.ParseInt(req.Form.Get("rows"), 10, 64)
		if err != nil {
			Error(rw, err)
			return
		}

		if err := wdb.Housekeeping(int(days), int(rows)); err != nil {
			Error(rw, err)
			return
		}
		web.Render().JSON(rw, http.StatusNoContent, nil)
	}
}
