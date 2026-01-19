import React from 'react'
import logoSvg from '../assets/logo.svg'
import './Logo.css'

const Logo = ({ size = 'default', showText = true }) => {
  const sizeClass = size === 'small' ? 'logo-small' : 'logo-default'
  
  return (
    <div className={`logo-container ${sizeClass}`}>
      <img src={logoSvg} alt="火影居家会" className="logo-svg" />
      {showText && <span className="logo-text">火影居家会</span>}
    </div>
  )
}

export default Logo
