function pair(name, value) {
    return name + "=" + encodeURIComponent(value)
}

class DefinedUrl {
    host = ""
    company = ""
    constructor(host, company) {
        this.host = host
        this.company = company
        
    }
    register() { 
        let url = `${this.host}/client/register`
        return url
    }
    login() { 
        let url =  `${this.host}/client/login` 
        return url
    }
    buyTickets(tocken, home, priceId) { 
        let url = `${this.host}/client/buy`
        url += `?${pair("Tocken", tocken)}`
        url += `&${pair("Home", home)}`
        url += `&${pair("PriceId", priceId)}`
        return url 
    }
    play(tocken, machId) { 
        let url = `${this.host}/machine/play`
        url += `?${pair("Tocken", tocken)}`
        url += `&${pair("MachId", machId)}` 
        return url
    }
    poll(tocken, eventId) { 
        let url = `${this.host}/machine/poll` 
        url += `?${pair("Tocken", tocken)}`
        url += `&${pair("EventId", eventId)}` 
        return url
    }
    pick(tocken, machId) {
        let url = `${this.host}/machine/getpub`
        url += `?${pair("Tocken", tocken)}`
        url += `&${pair("MachId", machId)}` 
        return url 
    }
    balance(tocken) { 
        let url = `${this.host}/client/balance` 
        url += `?${pair("Tocken", tocken)}`
        return url
    }
    prices(tocken) { 
        let url = `${this.host}/client/price` 
        url += `?${pair("Tocken", tocken)}`
        return url
    }
}

class Config {
    constructor(host, company) {
        this.host = host
        this.company = company
    }
    url(){
        return new DefinedUrl(this.host, this.company)
    }
   

}

export var config = new Config("https://clousel.fin-tech.com", "vesers")
// export var config = new Config("http://localhost:4321", "default")