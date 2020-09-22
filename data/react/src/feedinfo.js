import React, { Component } from 'react'
import TextField from '@material-ui/core/TextField'
import axios from "axios";
import { Button } from '@material-ui/core';
import Loading from './loading';

const env = process.env
console.log(process.env)
const apiUrl = env.REACT_APP_API
const FeedInfoApi = apiUrl + "/feedinfo"

export default class FeedInfo extends Component {
    constructor(props) {
        super(props);
        this.state = {
            feedInfoList: [],
            rss_url: "",
            xpath: "",
            loading: false,
        }
        this.handleSubmit = this.handleSubmit.bind(this);
        this.handleUrlChange = this.handleUrlChange.bind(this);
        this.handleXPathChange = this.handleXPathChange.bind(this);
    }

    componentDidMount() {
        this.setState( {loading: true}, async () => {
            try {
                const res = await axios.get(FeedInfoApi)
                this.setState({ feedInfoList: res.data })
            } catch (e) {
                alert("失敗しました" + e);
            }
            this.setState( {loading: false})
        })
    }

    onRSSClick = (id) => {
        window.open(apiUrl + "/rss/" + id, "_blank")
    }

    onDeleteClick(val) {
        var confirm = window.confirm(val.title + "を削除しますか?")
        if (confirm) {
            this.setState( {loading: true}, async () => {
                const deleteUrl = FeedInfoApi + "?id=" + val.id
                try {
                    await axios.delete(deleteUrl)
                    const res = await axios.get(FeedInfoApi)
                    this.setState({ feedInfoList: res.data })
                } catch (e) {
                    alert("失敗しました" + e);
                }
                this.setState( {loading: false})
            })
        }
    }

    handleUrlChange(event) {
        this.setState({ rss_url: event.target.value });
    }

    handleXPathChange(event) {
        this.setState({ xpath: event.target.value });
    }

    handleSubmit(event) {
        if (this.state.rss_url === '' || this.state.xpath === '') {
            alert('正しく入力してください')
            event.preventDefault();
            return
        }
        var res = window.confirm("この内容で登録しますか?\nRSS URL: " + this.state.rss_url + "\nXPath: " + this.state.xpath)
        if (res) {
            this.setState({ loading: true }, () => {
                axios.post(FeedInfoApi, {
                    "rss_url": this.state.rss_url,
                    "xpath": this.state.xpath
                }).then((_) => {
                    axios.get(FeedInfoApi).then((res) => {
                        console.log(res.data);
                        this.setState({ feedInfoList: res.data })
                        this.setState({ loading: false })

                    }).catch(e => {
                        alert("失敗しました" + e);
                        this.setState({ loading: false })
                    });
                }).catch(e => {
                    alert("失敗しました" + e);
                    this.setState({ loading: false })
                });
                this.setState({ rss_url: '', xpath: '' })
            })
            event.preventDefault();
        } else {
            event.preventDefault();
        }
    }

    render() {
        if (this.state.loading) {
            return (<Loading />);
        }
        return (<div align="center">
            <h1>Full RSS Feed Generator</h1>
            <form onSubmit={this.handleSubmit}>
                <label>
                    <TextField style={{ width: 300 }} id="rss-url" label="RSS URL" value={this.state.rss_url} onChange={this.handleUrlChange} />
                    <br />
                    <TextField style={{ width: 300 }} id="xpath" label="XPath" value={this.state.xpath} onChange={this.handleXPathChange} />
                </label>
                <br />
                <Button style={{ marginTop: 20 }} variant="contained" color="primary" type="submit">Register</Button>
            </form>

            {this.state.feedInfoList != null ? this.state.feedInfoList.map(val => (
                <ul key={val.id}>
                    <li style={{ listStyle: "none" }} align="center">
                        {val.title}
                        <Button style={{ padding: 10, margin: 5 }} variant="contained" disableElevation onClick={this.onRSSClick.bind(this, val.id)}>RSS</Button>
                        <Button style={{ padding: 10, margin: 5 }} variant="contained" disableElevation onClick={this.onDeleteClick.bind(this, val)}>Delete</Button>
                    </li>
                </ul>
            )) :
                <div></div>}
        </div>);
    }
}
