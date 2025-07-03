# `config` folder
In this config folder, there are some customizations that you can make to prevent run time overhead doing useless things.

## See the GoProject README.md for the main configurations!
Main configs:
- `UseSudo`
- `DoTracing`
- `CachedBuild`
- `SkipDialogue`

## Excluded services (`./excludeservices`)
In this file, one can specify which services to exclude from artifact generation furthermore (apart from default services defined in `../defaultdata/defaultservices`). So in case after the first run, you have the following service docker repos in `../artifacts/output/containers`:
- mongod
- irrelevantservice1
- irrelevantservice2

You can then define the following in `./excludeservices`:
`- irrelevantservice1
- irrelevantservice2`
