# Consumption-ms
It's and API that provides information regarding the energy consuption with three diferent filters or rates of time weekly. monthly and daily in a given window time.

# Run
You can find a file called Makefile to run it you will need to have "make" insatalled on your computer.

- you can be able to installed with the command below with 
linux:
`sudo apt install make`
macOS:
`brew install make`

__1.)__ Run without docker or air
`make run`
__2.)__ Run with air ( air is a Live Reload tool excellent for development )
`make run-air`

__3.)__ Run with docker
`make run-docker`
once you no longer need the container please run the following command
`make down-docker`

Once the project is up I suggest you go to the documentation built in swagger and make the request with the user
interface.

`http://localhost:8080/api/v1/docs/index.html`

###Example:
 Example to do the request
 `localhost:8080/api/v1/consumption?meter_ids=1,2&start_date=2023-05-30&end_date=2023-06-20&kind_period=weekly`
