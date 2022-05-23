import React from 'react';

const Home = React.lazy(() => import('./pages/Home/Home'));
const Vulnerabilities = React.lazy(() => import('./pages/Visualization/Vulnerabilities/Vulnerabilities'));
const List = React.lazy(() => import('./pages/Visualization/List/List'));

const routes = [
    {path: '/', exact: true, name: 'HFContractDetector', component: Home},
    {path: '/Visualization/Vulnerabilities', exact: true, name: '漏洞检测', component: Vulnerabilities},
    {path: '/Visualization/List', exact: true, name: '漏洞清单', component: List},
];

export default routes;
