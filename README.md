# gitlab2gs

-----------------

## Description

	migrate gitlab to gogs tool

## Installation

	$ go get github.com/WayneZhouChina/gitlab2gs

## Usage

	gitlab2gs -config=/path/config.json

### config file

	If "gitlabProjects" is empty, all projects will be migrated
	{
			"gitlabHost":     "https://gitlab.com",
			"gitlabApiPath": "/api/v3",
			"gitlabUser": "username",
			"gitlabPassword": "passwrod",
			"gitlabToken":    "xxxxx",
			"gitlabProjects": ["project1", "project2"],
			"gogsUrl": "gogs host url",
			"gogsToken": "yyyyyy",
			"gogsApiPath": "/api/v1"
	}

## AUTHOR

	Written by WayneZhou, cumtxhzyy[at]gmail.com

## COPYRIGHT

	Copyright (c) 2016 WayneZhou. This library is free software; you can redistribute it and/or modify it.

## Speical Thanks
	
	I have learnt a lot from [gitlab2gos](https://github.com/ewoutp/gitlab2gogs), Thanks to [ewoutp](https://github.com/ewoutp). I just add a field which is spcified projects you want.
