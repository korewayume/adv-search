import React, {useEffect, useState, useRef} from 'react';
import {useDynamicList, useDebounceFn} from 'ahooks';
import {Input, AutoComplete, Alert} from 'antd';
import lodash from 'lodash';
import axios from 'axios';

function AutoCompletedInput() {
  const input = useRef(null);
  const [inputValue, setInputValue] = useState(`ip="127.0.0.1" || 银行 && local && ! port="80" && (server="apache")`);
  const [inputCursor, setInputCursor] = useState(0);
  const [options, setOptions] = useState([]);
  const [dropdownOpen, setDropdownOpen] = useState(false);

  const errors = useDynamicList([]);

  const renderErrors = () => {
    return errors.list.map((error, index) => {
      return <Alert message={error} key={index} style={{"fontSize": "20px", "margin": "5px 0"}} type="error" closable/>
    })
  }

  const onChange = async (e) => {
    // 输入变化回调
    updateInput(e.target.value, e.target.selectionStart)
    await debounceSuggestions.run(e.target.value, e.target.selectionStart);
  }

  const selectionChangeHandle = async () => {
    // 光标变化回调
    updateInput(inputValue, input.current.input.selectionStart)
    await debounceSuggestions.run(inputValue, input.current.input.selectionStart);
  }

  useEffect(() => {
    document.addEventListener('selectionchange', selectionChangeHandle);
    return () => {
      document.removeEventListener('selectionchange', selectionChangeHandle);
    };
  });

  const setInputElementCursor = (position) => {
    // input.current.input.selectionStart === input.current.input.selectionEnd
    // 需要判断这个条件相等的原因是防止这个操作影响了用户的正常选中文本的操作
    if (input.current && input.current.input.selectionStart === input.current.input.selectionEnd) {
      input.current.input.selectionStart = position;
      input.current.input.selectionEnd = position
    }
  }

  useEffect(() => {
    setInputElementCursor(inputCursor)
  }, [input, inputCursor]);

  const onSelect = (value, option) => {
    setInputElementCursor(option.cursor)
    updateInput(value, option.cursor)
    setDropdownOpen(false)
  }

  const debounceSuggestions = useDebounceFn(
    async (input, cursor) => await requestSuggestions(input, cursor),
    {wait: 500}
  )

  const updateInput = (input, cursor) => {
    setInputValue(input);
    setInputCursor(cursor);
  }
  const requestSuggestions = async (input, cursor) => {
    try {
      const parse = await axios.post("/api/parse", {input})
      console.log(`POST /api/parse ${parse.status} (${parse.statusText})`, parse)
    } catch (error) {
      console.error(`POST /api/parse ${error.response.status} (${error.response.statusText})`, error.response)
    }
    try {
      const suggestion = await axios.post("/api/suggest", {cursor, input})
      if (suggestion.data.data) {
        setOptions(lodash.map(suggestion.data.data, item => {
          return {
            label: item.suggest,
            value: input.substring(0, item.start) + item.suggest + input.substring(item.end),
            cursor: item.start + item.suggest.length
          }
        }))
        setDropdownOpen(true)
      } else {
        errors.push(suggestion.data.error)
        if (errors.list.length > 4) {
          errors.shift()
        }
      }
    } catch (error) {
      if (error && error.response && error.response.data && error.response.data.error) {
        errors.push(error.response.data.error)
        if (errors.list.length > 4) {
          errors.shift()
        }
      }
    }
  }

  return (
    <div style={{"padding": "20px 10px", "fontSize": "34px"}}>
      <div style={{
        "display": "flex",
        "flexDirection": "column",
        "margin": "10px 0",
        "overflow": "auto"
      }}>
        {renderErrors()}
      </div>
      <AutoComplete
        style={{"width": "100%"}}
        options={options}
        onSelect={onSelect}
        defaultValue={inputValue}
        open={dropdownOpen}
        children={<Input ref={input} style={{"fontSize": "30px", "fontFamily": "Ubuntu Mono, monospaced"}} value={inputValue} onChange={onChange}/>}/>
    </div>
  )
}

export default AutoCompletedInput;
