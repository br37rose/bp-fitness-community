import React, { useState, useEffect } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

import { useRecoilState } from 'recoil';
import { currentUserState } from '../../../AppState';
import DatePicker from 'react-datepicker';
import 'react-datepicker/dist/react-datepicker.css';



const getSeriesName = (metricId, currentUser) => {
    if (metricId === currentUser.primaryHealthTrackingDeviceHeartRateMetricId) {
        return "Heart Rate";
    } else if (metricId === currentUser.primaryHealthTrackingDeviceStepsCountMetricId) {
        return "Steps Count";
    } else {
        return "Unknown Metric"; // default name if no match
    }
}



const ChartBuilder = ({ data }) => {
    const [currentUser] = useRecoilState(currentUserState);
    const [filterType, setFilterType] = useState('monthToDate'); // Default filter
    const [startDate, setStartDate] = useState(new Date());
    const [endDate, setEndDate] = useState(new Date());
    const [filteredData, setFilteredData] = useState([]);

    const applyDateFilter = (originalData, filterType, startDate, endDate) => {
        const now = new Date();
        let start, end;

        switch (filterType) {
            case 'monthToDate':
                start = new Date(now.getFullYear(), now.getMonth(), 1);
                end = now;
                break;
            case 'yearToDate':
                start = new Date(now.getFullYear(), 0, 1);
                end = now;
                break;
            case 'today':
                start = new Date(now.getFullYear(), now.getMonth(), now.getDate());
                end = new Date(now.getFullYear(), now.getMonth(), now.getDate() + 1);
                break;
            case 'custom':
                start = new Date(startDate);
                end = new Date(endDate);
                break;
            default:
                return originalData; // If no filter is applied, return the original data
        }

        return originalData.filter(entry => {
            const entryDate = new Date(entry.timestamp);
            return entryDate >= start && entryDate < end;
        });
    };


    const transformData = (originalData) => {
        // Group data by metricId
        const groupedData = originalData.reduce((acc, dataPoint) => {
            const metric = dataPoint.metricId;
            if (!acc[metric]) {
                acc[metric] = [];
            }
            acc[metric].push({
                category: new Date(dataPoint.timestamp).toLocaleTimeString(), // or any other format
                value: dataPoint.value
            });
            return acc;
        }, {});

        // Convert grouped data into the desired format
        const seriesData = Object.entries(groupedData).map(([metricId, data], index) => ({
            name: getSeriesName(metricId, currentUser),
            data
        }));


        return seriesData;
    };

    useEffect(() => {
        // Apply the filter first
        const dataWithFilter = applyDateFilter(data, filterType, startDate, endDate);

        // console.log(dataWithFilter)
        // Then transform the filtered data
        const transformedData = transformData(dataWithFilter);
        setFilteredData(transformedData);
    }, [data, filterType, startDate, endDate]);
    // console.log(filteredData)
    return (
        <div>
            <div>
                <label>
                    <input type="radio" value="monthToDate" checked={filterType === 'monthToDate'} onChange={() => setFilterType('monthToDate')} />
                    Month to Date
                </label>
                <label>
                    <input type="radio" value="yearToDate" checked={filterType === 'yearToDate'} onChange={() => setFilterType('yearToDate')} />
                    Year to Date
                </label>
                <label>
                    <input type="radio" value="today" checked={filterType === 'today'} onChange={() => setFilterType('today')} />
                    Just Today
                </label>
                <label>
                    <input type="radio" value="custom" checked={filterType === 'custom'} onChange={() => setFilterType('custom')} />
                    Custom Range
                </label>
                {filterType === 'custom' && (
                    <div>
                        <DatePicker selected={startDate} onChange={date => setStartDate(date)} />
                        <DatePicker selected={endDate} onChange={date => setEndDate(date)} />
                    </div>
                )}
            </div>
            <ResponsiveContainer width="100%" height={300}>
                <LineChart width={500} height={300}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="category" allowDuplicatedCategory={false} />
                    <YAxis dataKey="value" />
                    <Tooltip />
                    <Legend />
                    {filteredData.map((s) => (
                        <Line dataKey="value" data={s.data} name={s.name} key={s.name} />
                    ))}
                </LineChart>
            </ResponsiveContainer>
        </div>

    );
};

export default ChartBuilder;


