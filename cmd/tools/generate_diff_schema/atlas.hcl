data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/models/entity",
    "--dialect", "postgres", // | postgres | sqlite | sqlserver
  ]
}
env "local" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/15/dev" # dev database
  migration {
    dir = "file://./migrations"
    format = golang-migrate
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }

}

lint {
  destructive {
    force = false
  }
}

