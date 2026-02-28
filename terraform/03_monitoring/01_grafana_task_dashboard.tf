resource "grafana_dashboard" "chronflow_task_grafana_dashboard" {
  config_json = jsonencode({
    "annotations" : {
      "list" : [
        {
          "builtIn" : 1,
          "datasource" : {
            "type" : "grafana",
            "uid" : "-- Grafana --"
          },
          "enable" : true,
          "hide" : true,
          "iconColor" : "rgba(0, 211, 255, 1)",
          "name" : "Annotations & Alerts",
          "type" : "dashboard"
        }
      ]
    },
    "editable" : true,
    "fiscalYearStartMonth" : 0,
    "graphTooltip" : 0,
    "links" : [],
    "panels" : [
      {
        "datasource" : {
          "type" : "prometheus",
          "uid" : "prometheus"
        },
        "fieldConfig" : {
          "defaults" : {
            "color" : {
              "mode" : "thresholds"
            },
            "mappings" : [],
            "thresholds" : {
              "mode" : "absolute",
              "steps" : [
                {
                  "color" : "green",
                  "value" : 0
                },
                {
                  "color" : "red",
                  "value" : 80
                }
              ]
            }
          },
          "overrides" : []
        },
        "gridPos" : {
          "h" : 8,
          "w" : 12,
          "x" : 0,
          "y" : 0
        },
        "id" : 1,
        "options" : {
          "colorMode" : "value",
          "graphMode" : "area",
          "justifyMode" : "auto",
          "orientation" : "auto",
          "percentChangeColorMode" : "inverted",
          "reduceOptions" : {
            "calcs" : [
              "lastNotNull"
            ],
            "fields" : "",
            "values" : false
          },
          "showPercentChange" : true,
          "textMode" : "value_and_name",
          "wideLayout" : true
        },
        "pluginVersion" : "12.4.0",
        "targets" : [
          {
            "datasource" : {
              "type" : "prometheus",
              "uid" : "prometheus"
            },
            "editorMode" : "builder",
            "expr" : "chronflow_api_task_processed{status=\"success\"}",
            "legendFormat" : "Success",
            "range" : true,
            "refId" : "A"
          },
          {
            "datasource" : {
              "type" : "prometheus",
              "uid" : "prometheus"
            },
            "editorMode" : "builder",
            "expr" : "chronflow_api_task_processed{status=\"fail\"}",
            "instant" : false,
            "legendFormat" : "Fail",
            "range" : true,
            "refId" : "B"
          }
        ],
        "title" : "[API] Task processed",
        "type" : "stat"
      },
      {
        "datasource" : {
          "type" : "prometheus",
          "uid" : "prometheus"
        },
        "fieldConfig" : {
          "defaults" : {
            "color" : {
              "mode" : "thresholds"
            },
            "mappings" : [],
            "thresholds" : {
              "mode" : "absolute",
              "steps" : [
                {
                  "color" : "green",
                  "value" : 0
                },
                {
                  "color" : "red",
                  "value" : 80
                }
              ]
            }
          },
          "overrides" : []
        },
        "gridPos" : {
          "h" : 8,
          "w" : 12,
          "x" : 12,
          "y" : 0
        },
        "id" : 2,
        "options" : {
          "colorMode" : "value",
          "graphMode" : "area",
          "justifyMode" : "auto",
          "orientation" : "auto",
          "percentChangeColorMode" : "inverted",
          "reduceOptions" : {
            "calcs" : [
              "lastNotNull"
            ],
            "fields" : "",
            "values" : false
          },
          "showPercentChange" : true,
          "textMode" : "auto",
          "wideLayout" : true
        },
        "pluginVersion" : "12.4.0",
        "targets" : [
          {
            "editorMode" : "builder",
            "expr" : "chronflow_scheduler_task_scheduled",
            "legendFormat" : "__auto",
            "range" : true,
            "refId" : "A"
          }
        ],
        "title" : "[Scheduler] Task scheduled",
        "type" : "stat"
      },
      {
        "datasource" : {
          "type" : "prometheus",
          "uid" : "prometheus"
        },
        "fieldConfig" : {
          "defaults" : {
            "color" : {
              "mode" : "thresholds"
            },
            "mappings" : [],
            "thresholds" : {
              "mode" : "absolute",
              "steps" : [
                {
                  "color" : "green",
                  "value" : 0
                },
                {
                  "color" : "red",
                  "value" : 80
                }
              ]
            }
          },
          "overrides" : []
        },
        "gridPos" : {
          "h" : 8,
          "w" : 12,
          "x" : 0,
          "y" : 8
        },
        "id" : 3,
        "options" : {
          "colorMode" : "value",
          "graphMode" : "area",
          "justifyMode" : "auto",
          "orientation" : "auto",
          "percentChangeColorMode" : "inverted",
          "reduceOptions" : {
            "calcs" : [
              "lastNotNull"
            ],
            "fields" : "",
            "values" : false
          },
          "showPercentChange" : true,
          "textMode" : "value_and_name",
          "wideLayout" : true
        },
        "pluginVersion" : "12.4.0",
        "targets" : [
          {
            "editorMode" : "builder",
            "expr" : "chronflow_worker_task_processed{status=\"success\"}",
            "legendFormat" : "Success",
            "range" : true,
            "refId" : "A"
          },
          {
            "datasource" : {
              "type" : "prometheus",
              "uid" : "prometheus"
            },
            "editorMode" : "builder",
            "expr" : "chronflow_worker_task_processed{status=\"fail\"}",
            "instant" : false,
            "legendFormat" : "Fail",
            "range" : true,
            "refId" : "B"
          }
        ],
        "title" : "[Worker] Task status",
        "type" : "stat"
      },
      {
        "datasource" : {
          "type" : "prometheus",
          "uid" : "prometheus"
        },
        "fieldConfig" : {
          "defaults" : {
            "color" : {
              "mode" : "thresholds"
            },
            "mappings" : [],
            "thresholds" : {
              "mode" : "absolute",
              "steps" : [
                {
                  "color" : "green",
                  "value" : 0
                },
                {
                  "color" : "red",
                  "value" : 80
                }
              ]
            }
          },
          "overrides" : []
        },
        "gridPos" : {
          "h" : 8,
          "w" : 12,
          "x" : 12,
          "y" : 8
        },
        "id" : 4,
        "options" : {
          "colorMode" : "value",
          "graphMode" : "area",
          "justifyMode" : "auto",
          "orientation" : "auto",
          "percentChangeColorMode" : "inverted",
          "reduceOptions" : {
            "calcs" : [
              "lastNotNull"
            ],
            "fields" : "",
            "values" : false
          },
          "showPercentChange" : true,
          "textMode" : "auto",
          "wideLayout" : true
        },
        "pluginVersion" : "12.4.0",
        "targets" : [
          {
            "editorMode" : "builder",
            "expr" : "chronflow_janitor_task_recovered",
            "legendFormat" : "__auto",
            "range" : true,
            "refId" : "A"
          }
        ],
        "title" : "[Janitor] Tasks recovered",
        "type" : "stat"
      }
    ],
    "preload" : false,
    "refresh" : "5s",
    "schemaVersion" : 42,
    "tags" : [],
    "templating" : {
      "list" : []
    },
    "time" : {
      "from" : "now-1h",
      "to" : "now"
    },
    "timepicker" : {},
    "timezone" : "browser",
    "title" : "[Chronflow] Task",
    "uid" : "adtgvbn",
    "version" : 3,
    "weekStart" : ""
  })
}
