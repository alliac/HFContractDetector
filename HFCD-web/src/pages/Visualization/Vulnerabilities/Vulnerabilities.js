import 'antd/dist/antd.css'
import '../style.css'
import React, {useState} from 'react'
import {Alert, Button, Card, Col, Form, Input, message, Row, Select, Spin, Typography, Upload} from 'antd'
import {InboxOutlined} from '@ant-design/icons'
import {Prism} from 'react-syntax-highlighter'
import {atomDark} from 'react-syntax-highlighter/dist/esm/styles/prism'


const {Dragger} = Upload;
const {Option} = Select;
const {TextArea} = Input;
const {Text} = Typography
const {Paragraph} = Typography;
const Vulnerabilities = () => {

    const [code, setCode] = useState('\n\n\n\n\n\n\n\n\n\n')
    const [result, setResult] = useState('')
    const [loading, setLoading] = useState(false)
    const [form] = Form.useForm();
    let newcode;
    const ChaincodeDeal = (chaincode) => {
        newcode = code.replaceAll(chaincode, "ChainCode")
        newcode = newcode.replaceAll(/package .+/g, "package models")
    }
    const onFinish = (values) => {
        console.log('Success:', values);
        setResult('');
        setLoading(true)
        let data;
        let chaincode;
        //let newcode;
        let url;
        chaincode = form.getFieldValue("chaincode")
        ChaincodeDeal(chaincode)
        try {
            data = {
                "version": values.version,
                "newcode": newcode,
            }
            if (values.detect === 0) {
                url = "http://192.168.93.141:8080/"
                //url = "http://192.168.93.141:8080/static_detect"
            } else if (values.detect === 1) {
                url = "http://192.168.93.141:8080/dynamic_detect"
            } else {
                url = "http://192.168.93.141:8080/mix_detect"
            }
            console.log(data)
            const axios = require('axios');
            axios.get(url, {
                headers: {
                    'Content-type': 'multipart/form-data',
                    'Access-Control-Allow-Origin': '*',
                    'Access-Control-Allow-Headers': 'Content-Type, Authorization',
                },
                credentials: 'same-origin',
                withCredentials: false,
                timeout: 0,
                proxy: {
                    host: 'localhost',
                    port: 8080
                },
                ServerTimeOut: 1000000,
                params: data
            }).then(function (response) {
                setLoading(false)
                console.log(response.data["Result"])
                setResult(response.data["Result"])
            }).catch(function (error) {
                setLoading(false)
                if (error.response) {
                    console.log(error.response.headers);
                } else {
                    setResult(error.message)
                    console.log(error.message);
                }
                console.log(error.config);
            });
        } catch (err) {
            alert(err)
        }
    };

    const onFinishFailed = (errorInfo) => {
        console.log('Failed:', errorInfo);
    };

    const props = {
        name: 'file',
        multiple: true,
        action: '',
        onChange(info) {
            const {status} = info.file;
            if (status !== 'uploading') {
                console.log(info.file, info.fileList);
            }
            if (status === 'done') {
                message.success(`${info.file.name} file uploaded successfully.`);
            } else if (status === 'error') {
                message.error(`${info.file.name} file upload failed.`);
            }
            if (status === 'removed') {
                setCode('\n\n\n\n\n\n\n\n\n\n');
                setResult('');
            }
        },
        onDrop(e) {
            console.log('Dropped files', e.dataTransfer.files);
        },
    };
    const prism = (
        <div className="example" style={{height: 275, overflow: 'auto'}}>
            <Prism className="example"
                   showLineNumbers={true}
                   startingLineNumber={1}
                   language="go" style={atomDark}
                   wrapLines={true}
                   lineNumberStyle={{color: '#fff', fontSize: 5}}>
                {code}
            </Prism>
        </div>

    )
    const container = (
        <Alert className="example"
               message={prism}
            //description="34"
               type="info"
        />
    );

    return (
        <>
            <div className="site-card-border-less-wrapper">
                <Row gutter={3}>
                    <Col span={5}>
                        <Card bordered={false} style={{height: 550}}>
                            <Dragger
                                beforeUpload={file => {
                                    const reader = new FileReader()
                                    reader.readAsText(file);
                                    reader.addEventListener('load', event => {
                                        const txt = event.target.result;
                                        //console.log(txt)
                                        setCode(txt)
                                    })
                                    return false
                                }}
                                {...props}>
                                <p className="ant-upload-drag-icon">
                                    <InboxOutlined/>
                                </p>
                                <p className="ant-upload-text">上传文件</p>
                            </Dragger>
                        </Card>
                    </Col>
                    <Col span={14}>
                        <Card>
                            <Text code>智能合约</Text>
                            <Spin spinning={loading} delay={200} tip={"正在检测"}>
                                {container}
                            </Spin>
                        </Card>
                        <Card style={{height: 200}}>
                            <div>
                                <Paragraph copyable={{text: result}}><Text code>测试结果</Text></Paragraph>
                                <TextArea bordered={false}
                                          value={result}
                                          placeholder=""
                                          autoSize={{minRows: 6, maxRows: 6}}>
                                </TextArea>
                            </div>
                        </Card>
                    </Col>
                    <Col span={5}>
                        <Card bordered={false} style={{height: 550}}>
                            <Form
                                form={form}
                                name="basic"
                                onFinish={onFinish}
                                onFinishFailed={onFinishFailed}
                                autoComplete="off">
                                <Form.Item
                                    name="chaincode"
                                    rules={[
                                        {
                                            required: true,
                                            message: 'Please input chaincode!',
                                        },
                                    ]}>
                                    <Input placeholder="请输入Chaincode对象名称" name={"chaincode"}/>
                                </Form.Item>
                                <Form.Item
                                    name="version"
                                    rules={[
                                        {
                                            required: true,
                                            message: 'Please input fabric version!',
                                        },
                                    ]}>
                                    <Select
                                        placeholder="请选择Fabric版本"
                                    >
                                        <Option value="1.1">1.1</Option>
                                        <Option value="1.2">1.2</Option>
                                        <Option value="1.4">1.4</Option>
                                    </Select>
                                </Form.Item>
                                <Form.Item
                                    name="detect"
                                    rules={[
                                        {
                                            required: true,
                                            message: 'Please input detect method!',
                                        },
                                    ]}>
                                    <Select
                                        placeholder="请选择检测方式"

                                        allowClear>
                                        <Option value={0}>静态检测</Option>
                                        <Option value={1}>动态检测</Option>
                                        <Option value={2}>混合检测</Option>
                                    </Select>
                                </Form.Item>
                                <Form.Item style={{textAlign: 'center'}}>
                                    <Button type="primary" htmlType="submit">
                                        开始检测
                                    </Button>
                                </Form.Item>
                            </Form>
                        </Card>
                    </Col>
                </Row>
            </div>
        </>
    );
};
export default Vulnerabilities