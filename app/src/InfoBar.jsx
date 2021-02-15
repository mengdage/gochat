import FlexSpacer from './helper/FlexSpacer'
import './InfoBar.css'

const InfoBar = ({toUser, serverAddr}) => {
    return (
        <div className="infobar-container">
            <div className="infobar-touser">TO: {toUser}</div>
            <FlexSpacer />
            <div></div>
        </div>
    )
}

export default InfoBar
