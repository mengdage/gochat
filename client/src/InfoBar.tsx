import FlexSpacer from './helper/FlexSpacer'
import './InfoBar.scss'

interface InfoBarProps {
    toUserName: string
    serverAddr: string
}
const InfoBar: React.FC<InfoBarProps> = ({toUserName, serverAddr}) => {
    return (
        <div className="infobar-container">
            <div className="infobar-touser">TO {toUserName}</div>
            <FlexSpacer />
            <div className="infobar-extra-info">{serverAddr}</div>
        </div>
    )
}

export default InfoBar
