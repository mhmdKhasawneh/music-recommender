import React, { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
function Main() {
    const navigate = useNavigate()
    const [recommendTo, setRecommendTo] = useState("")
    const [url, setUrl] = useState("")
    const [recommendations, setRecommendations] = useState([])

    useEffect(() => {
        checkAuthStatus()
        displayReccommendations()
    }, [])

    const checkAuthStatus = async () => {
        const response = await fetch("http://localhost:8080/api/user/tokenlogin", {
            mode: "cors",
            method: "GET",
            headers: {
                'Authorization': "Bearer ".concat(localStorage.getItem("Authorization"))
            },
        })
        if (!response.ok) {
            navigate("/login")
        }
    }
    const displayReccommendations = async () => {
        const response = await fetch("http://localhost:8080/api/recommendation/get", {
            mode: "cors",
            method: "GET",
            headers: {
                'Authorization': "Bearer ".concat(localStorage.getItem("Authorization"))
            },
        })
        const data = await response.json()
        setRecommendations(data["recommendations"])
    }

    const recommend = async () => {
        const response = await fetch("http://localhost:8080/api/recommendation/recommend", {
            mode: "cors",
            method: "POST",
            body: JSON.stringify({
                'to_user': recommendTo,
                'url': url
            }),
            headers: {
                'Content-Type': 'application/json',
                'Authorization': "Bearer ".concat(localStorage.getItem("Authorization"))
            },
        })
        if (!response.ok) {
            if (response.status == 401) {
                navigate('/login')
            }
            else {
                window.alert("Status code " + response.status + "\nMake sure you entered a valid URL or email")
            }
        }

    }

    return (
        <div>
            <div className="container-before"></div>
            <div className='recommendation'>
                <input className="input" type="text" placeholder='Enter URL here. Only Spotify is accepted'
                    onChange={(e) => setUrl(e.target.value)} />
                <input className="input" type="text" placeholder="Enter the user's email you want to recommend to here"
                    onChange={(e) => setRecommendTo(e.target.value)} />
                <br />
                <button onClick={recommend}>RECOMMEND</button>
            </div>
            <div className="container-after"></div>
            {
            recommendations.length == 0 
                ?
                <h1>LOSER GOT NO RECOMMENDATIONS</h1>
                :
                <ul>
                    {recommendations.map(recommendation => {
                        return <li key={recommendation['ID']}>
                            <img src={recommendation['ImgUrl']} alt="image" />
                            <span className="track">{recommendation['Name']}</span>
                            <a className='nostyle' href={recommendation['Url']}>
                                <span className="redirect">Listen on Spotify</span>
                            </a>
                        </li>
                    })}
                </ul>
            }
        </div>
    )
}

export default Main