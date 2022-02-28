import React from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'


const Instructions = (props) => {
    const markdown = window.atob(props.yaml);
    return (
        <ReactMarkdown children={markdown} remarkPlugins={[remarkGfm]} />
    );
}
export default Instructions;