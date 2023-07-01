package main

var (
	PLAYBOOKS = map[string]func(){
		"deploy": playbookDeploy,
		"start":  playbookStart,
		"stop":   playbookStop,
		"reset":  playbookReset,
		"stats":  playbookStats,
		"exec":   playbookExec,
	}
)
